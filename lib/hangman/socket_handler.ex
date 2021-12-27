defmodule Hangman.SocketHandler do
    @behaviour :cowboy_websocket

    @lives 10
    @guess_delay 5
    @ping_interval 5_000

    def init(request, _state) do
        {:cowboy_websocket, request, %{}, %{
            idle_timeout: @ping_interval * 2
        }}
    end
    def websocket_init(state), do: {:ok, state}
    def websocket_handle({:text, json}, state) do
        case Jason.decode(json) do
            {:ok, payload} ->
                case payload["action"] do
                    "create_game" ->
                        if valid_game_type(payload["game_type"]) && !Map.has_key?(state, "game_id") do
                            game_id = :crypto.strong_rand_bytes(8) |> Base.encode16
                            word = Hangman.RandWord.rand_word()
                            game_type = String.to_atom(payload["game_type"])
                            :ets.insert(:game_data, {game_id, game_type, word, @lives, [], 0, :a, false})
                            Registry.Hangman |>
                                Registry.register(game_id, {})
                            if game_type === :local do
                                fan_out(game_id, Jason.encode!(
                                    %{
                                        message: "update",
                                        word: hide_word(word, []),
                                        letters: [],
                                        lives: @lives,
                                        player: :a,
                                        full: true,
                                        guess_word_locked: false
                                    }
                                ), true)
                                {:ok, %{"game_id" => game_id}}
                            else
                                {:reply, {:text, Jason.encode!(
                                    %{
                                        message: "wait_for_partner",
                                        game_id: game_id
                                    }
                                )}, %{"game_id" => game_id}}
                            end
                        else
                            {:stop, state}
                        end
                    "join_game" ->
                        game_id = payload["game_id"]
                        game_id_valid = is_valid_game_id(game_id)
                        already_in_game_or_local_game = if game_id_valid do
                            Map.has_key?(state, "game_id") || game_type(game_id) === :local
                        else
                            false
                        end
                        if already_in_game_or_local_game do
                            {:stop, state}
                        else
                            if !game_id_valid || game_full(game_id) do
                                {:reply, {:text, Jason.encode!(
                                    %{
                                        message: "cannot_join_game"
                                    }
                                )}, %{}}
                            else
                                Registry.Hangman |>
                                    Registry.register(game_id, {})
                                :ets.insert(:game_data, {
                                    game_id,
                                    game_type(game_id),
                                    game_word(game_id),
                                    game_lives(game_id),
                                    used_letters(game_id),
                                    last_guess(game_id),
                                    cur_player(game_id),
                                    true
                                })
                                fan_out(game_id, Jason.encode!(
                                    %{
                                        message: "update",
                                        word: hide_word(game_word(game_id), []),
                                        letters: [],
                                        lives: @lives,
                                        player: :a,
                                        full: true,
                                        guess_word_locked: false
                                    }
                                ), true)
                                {:ok, %{"game_id" => game_id}}
                            end
                        end
                    "guess_letter" ->
                        game_id = state["game_id"]
                        if !is_valid_game_id(game_id) do
                            {:stop, state}
                        else
                            if game_lives(game_id) <= 0 || Enum.member?(
                                used_letters(game_id),
                                payload["letter"]
                            ) do
                                {:stop, state}
                            else
                                lives = if String.contains?(game_word(game_id), payload["letter"]) do
                                    game_lives(game_id)
                                else
                                    game_lives(game_id) - 1
                                end
                                letters = if lives > 0 do
                                    [payload["letter"] | used_letters(game_id)]
                                else
                                    game_word(game_id) |> String.split("", trim: true)
                                end
                                player = swap_player(game_id)
                                :ets.insert(:game_data, {
                                    game_id,
                                    game_type(game_id),
                                    game_word(game_id),
                                    lives,
                                    letters,
                                    last_guess(game_id),
                                    player,
                                    true
                                })
                                fan_out(game_id, Jason.encode!(
                                    %{
                                        message: "update",
                                        word: hide_word(game_word(game_id), letters),
                                        letters: letters,
                                        lives: lives,
                                        player: player,
                                        full: true,
                                        guess_word_locked: !can_guess(game_id)
                                    }
                                ), true)
                                {:ok, state}
                            end
                        end
                    "guess_word" ->
                        game_id = state["game_id"]
                        if !is_valid_game_id(game_id) do
                            {:stop, state}
                        else
                            if can_guess(game_id) do
                                if payload["word"] === game_word(game_id) do
                                    :ets.insert(:game_data, {
                                        game_id,
                                        game_type(game_id),
                                        game_word(game_id),
                                        game_lives(game_id),
                                        game_word(game_id) |> String.split("", trim: true),
                                        System.system_time(:second),
                                        swap_player(game_id),
                                        true
                                    })
                                else
                                    fan_out(game_id, Jason.encode!(
                                        %{
                                            message: "lock_guess_word"
                                        }
                                    ), false)
                                    :ets.insert(:game_data, {
                                        game_id,
                                        game_type(game_id),
                                        game_word(game_id),
                                        game_lives(game_id),
                                        used_letters(game_id),
                                        System.system_time(:second),
                                        swap_player(game_id),
                                        true
                                    })
                                end
                                fan_out(game_id, Jason.encode!(
                                    %{
                                        message: "update",
                                        word: hide_word(game_word(game_id), used_letters(game_id)),
                                        letters: used_letters(game_id),
                                        lives: game_lives(game_id),
                                        player: cur_player(game_id),
                                        full: true,
                                        guess_word_locked: true
                                    }
                                ), true)
                            end
                            {:ok, state}
                        end
                    "new_game" ->
                        game_id = state["game_id"]
                        if !is_valid_game_id(game_id) do
                            {:stop, state}
                        else
                            hidden_word = hide_word(game_word(game_id), used_letters(game_id))
                            if game_lives(game_id) <= 0 || !String.contains?(hidden_word, "_") do
                                word = Hangman.RandWord.rand_word()
                                :ets.insert(:game_data, {
                                    game_id,
                                    game_type(game_id),
                                    word,
                                    @lives,
                                    [],
                                    0,
                                    :a,
                                    true
                                })
                                fan_out(game_id, Jason.encode!(
                                    %{
                                        message: "update",
                                        word: hide_word(word, []),
                                        letters: [],
                                        lives: @lives,
                                        player: :a,
                                        full: true,
                                        guess_word_locked: false
                                    }
                                ), true)
                                {:ok, state}
                            else
                                {:stop, state}
                            end
                        end
                    "ping" ->
                        {:ok, state}
                    _ ->
                        {:stop, state}
                end
            {:error, _} -> {:stop, state}
        end
    end
    def websocket_info(info, state), do: {:reply, {:text, info}, state}
    def terminate(_, _, state) do
        if is_valid_game_id(state["game_id"]) do
            fan_out(state["game_id"], Jason.encode!(
                %{
                    message: "abandoned"
                }
            ), false)
            :ets.delete(:game_data, state["game_id"])
        end
        :ok
    end
    defp is_valid_game_id(game_id), do: length(:ets.lookup(:game_data, game_id)) > 0
    defp lookup_game_data(game_id), do: :ets.lookup(:game_data, game_id) |> Enum.at(0)
    defp game_type(game_id), do: lookup_game_data(game_id) |> elem(1)
    defp game_word(game_id), do: lookup_game_data(game_id) |> elem(2)
    defp game_lives(game_id), do: lookup_game_data(game_id) |> elem(3)
    defp used_letters(game_id), do: lookup_game_data(game_id) |> elem(4)
    defp last_guess(game_id), do: lookup_game_data(game_id) |> elem(5)
    defp cur_player(game_id), do: lookup_game_data(game_id) |> elem(6)
    defp game_full(game_id), do: lookup_game_data(game_id) |> elem(7)
    defp can_guess(game_id), do: System.system_time(:second) - last_guess(game_id) > @guess_delay
    defp valid_game_type(game_type), do: game_type === "local" || game_type === "remote"
    defp swap_player(game_id), do: if cur_player(game_id) === :a, do: :b, else: :a
    defp hide_word(word, letters) do
        word
            |> String.split("", trim: true)
            |> Enum.map(&(if Enum.member?(letters, &1), do: &1, else: "_"))
            |> Enum.join("")
    end
    defp fan_out(game_id, message, include_self) do
        Registry.Hangman |>
            Registry.dispatch(game_id, fn entries ->
                for {pid, _} <- entries do
                    if pid != self() || include_self do
                        Process.send(pid, message, [])
                    end
                end
            end)
    end
end
