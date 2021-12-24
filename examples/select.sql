SELECT a, b, AVG( (a+b)/c*d) ),c, SUM(a-b) AS total FROM x INNER JOIN y ON x.a = y.b GROUP BY a,b ORDER BY a, b DESC
