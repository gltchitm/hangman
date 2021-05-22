defmodule Hangman.StaticResources do
    use Plug.Builder

    plug(
        Plug.Static,
        at: "/",
        from: :hangman
    )
    plug(:not_found)

    def not_found(conn, _) do
        send_resp(conn, 404, "Not Found")
    end
end
