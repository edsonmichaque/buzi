{{ range $k, $v := .Types }}
class {{$k}}:
    pass
{{ end }}
