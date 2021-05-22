defmodule Hangman.Application do
    use Application

    def start(_, _) do
        children = [
            Plug.Cowboy.child_spec(
                scheme: :http,
                plug: Hangman.Router,
                options: [
                    port: System.get_env("PORT") || 5522,
                    dispatch: dispatch()
                ]
            ),
            Registry.child_spec(
                keys: :duplicate,
                name: Registry.Hangman
            )
        ]
        options = [strategy: :one_for_one, name: Hangman.Supervisor]

        :ets.new(:game_data, [:set, :public, :named_table])

        Supervisor.start_link(children, options)
    end
    defp dispatch() do
        [
            {:_,
                [
                    {"/ws", Hangman.SocketHandler, []},
                    {:_, Plug.Cowboy.Handler, {Hangman.Router, []}}
                ]
            }
        ]
    end
end
