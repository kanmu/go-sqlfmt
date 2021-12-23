SELECT a, b, AVG( (a+b)/c*d) ) FROM x INNER JOIN y ON x.a = y.b GROUP BY a,b ORDER BY a, b DESC
