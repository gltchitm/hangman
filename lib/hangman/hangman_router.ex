defmodule Hangman.Router do
    use Plug.Router

    plug(:match)
    plug(:dispatch)

    forward("/static", to: Hangman.StaticResources)

    get "/" do
        send_resp(conn, 200, EEx.eval_file("lib/templates/index.html.eex"))
    end

    match _ do
        send_resp(conn, 404, "Not Found")
    end
end
