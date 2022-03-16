SELECT * FROM product
WHERE name in (
    {{range $i, $p := . -}}
    {{- if ne $i 0 }}
    ,
    {{- end }}
    :name{{- $i}}
    {{ end -}}
)
;
