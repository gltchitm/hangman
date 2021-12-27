defmodule Hangman.MixProject do
    use Mix.Project

    def project do
        [
            app: :hangman,
            version: "1.0.0",
            elixir: "~> 1.12",
            start_permanent: Mix.env() === :prod,
            deps: deps()
        ]
    end
    def application do
        [
            mod: {Hangman.Application, []}
        ]
    end
    defp deps do
        [
            {:plug_cowboy, "~> 2.0"},
            {:jason, "~> 1.2"}
        ]
    end
end
