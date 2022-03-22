SELECT *, COUNT(*) OVER()
{{ range $i, $q := .Q -}}
{{- if eq $i 0 -}}
FROM
{{- else -}}
INNER JOIN
{{- end }}
(
{{- range $j, $member := $.TargetFields -}}
    {{ if ne $j 0 }}
    UNION
    {{- end }}
    SELECT * FROM product
    {{- if eq $.PatternMatch "exact"}}
    WHERE {{ $member }} = :q{{- $i}}
    {{ else }}
    WHERE {{ $member }} LIKE concat('%', cast(:q{{- $i}} as text), '%')
    {{- end }}
{{- end }}
) as T{{- $i }}
{{ if ne $i 0}}on T0.id = T{{- $i }}.id{{ end }}
{{ end -}}
ORDER BY T0.id LIMIT :limit OFFSET :offset
