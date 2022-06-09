module {{.Params.module}}
    {{ range $k, $v := .Types }}
    class {{$k}}
    end
    {{ end }}
end