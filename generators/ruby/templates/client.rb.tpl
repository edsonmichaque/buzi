module {{.Params.module}}
    class Client
    {{ range $k, $v := .Operations }}
        def {{snakecase $k}}()
            return nil
        end
    {{ end }}
    end
end