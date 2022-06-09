package {{.Params.package}}

{{ range $k, $v := .Types }}
type {{$k}} struct {
    {{ if isobj $v -}}
        {{ range $name, $kind := props $v -}}
            {{ if isscalar $kind -}}
                {{camelcase $name }} {{ kind $kind }} `json:"{{$name}},omitempty"`
            {{ end -}}


            {{ if isarray $kind -}}
                {{ if isscalar $kind.Value.Items }}
                    {{camelcase $name }} []{{ $kind.Value.Items.Value.Type }} `json:"{{$name}},omitempty"`
                {{ end }}
            {{ end -}}
        {{ end -}}

        {{ range $name, $kind := refs $v -}}
            {{ if isref $kind -}}
                {{camelcase $name }} *{{ ref $kind }} `json:"{{$name}},omitempty"`
            {{ end -}}
        {{ end -}}

    {{ end }}
}
{{ end }}