class Client:
{{ range $k, $v := .Operations }}
    def {{snakecase $k}}(self):
        pass
{{ end }}
