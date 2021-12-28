WITH ls AS (select n,ARRAY[x,y] as 'listed', CASE when x>0 then true else false END as flagged from lists),
v AS (select first_Value(x) OVER(z order by u,v) as fv FROM values)
SELECT CASE a<b THEN a ELSE b END as minab, a, b, AVG( (a+b)/c*d) ),c, user, user() AS user_name, current_user() AS user_name2,
to_Char(date_t) AS 'mydate', to_date('date_str')::timestamp AS 'my_other_date',
left(a, 5) AS truncated_a, ROW_number() as row_num,
jsonb_agg(json_build_object(key, value)::jsonb)->>'key', E'abc' AS literal, U&'\0x' AS utf,
ls.'listed', v.fv, CASE WHEN LOG10('listed'[0]) IS NULL THEN 'null'
WHEN LOG10('listed'[0]) = NaN THEN 'wrong' ELSE 'ok' END, ST_POINT(x)::geography,
SUM(a-b) AS total FROM x INNER JOIN y ON x.a = y.b LEFT JOIN ls ON ls.n = y.c CROSS Join "v" GROUP BY a,b ORDER BY a, b DESC
