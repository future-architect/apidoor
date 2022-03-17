SELECT pd.id, pd.contract_id, pd.product_id FROM
(
    {{ range $i, $c := . -}}
    {{- if ne $i 0 -}} UNION ALL {{- end }}
    SELECT id, contract_id, product_id
    FROM contract_product_content
    WHERE contract_id = :contract_id_{{- $i}}
        {{- if $c.ProductIDs }}
        AND product_id IN
            ( {{ range $j, $p := $c.ProductIDs -}}
                    {{- if ne $j 0 -}},{{- end }}:product_id_{{- $i -}}_{{- $j -}}
            {{- end }} )
        {{- end }}
    {{ end }}
) as pd INNER JOIN (
    SELECT id FROM contract
    WHERE user_id = :user_id
) as ct on pd.contract_id = ct.id;


