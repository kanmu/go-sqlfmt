select a, b, aVg( (a+b)/c*d) ) from x left outer join  y on x.a = y.b group by a,b order by a, b desc limit 5
