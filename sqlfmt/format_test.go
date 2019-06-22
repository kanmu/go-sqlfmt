package sqlfmt

import (
	"testing"
)

func TestCompare(t *testing.T) {
	test := struct {
		before string
		after  string
		want   bool
	}{
		before: "select * from xxx",
		after:  "select\n  *\nFROM xxx",
		want:   true,
	}
	if got := compare(test.before, test.after); got != test.want {
		t.Errorf("want %v#v got %#v", test.want, got)
	}
}

func TestRemove(t *testing.T) {
	got := removeSpace("select xxx from xxx")
	want := "selectxxxfromxxx"
	if got != want {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

func TestFormat(t *testing.T) {
	for _, tt := range formatTestingData {
		opt := &Options{}
		t.Run(tt.src, func(t *testing.T) {
			got, err := Format(tt.src, opt)
			if err != nil {
				t.Errorf("should be nil, got %v", err)
			}
			if tt.want != got {
				t.Errorf("\nwant %#v, \ngot %#v", tt.want, got)
			}
		})
	}
}

var formatTestingData = []struct {
	src  string
	want string
}{
	{
		src: `select name, age, id from user join transaction on a = b`,
		want: `
SELECT
  name
  , age
  , id
FROM user
JOIN transaction
ON a = b`,
	},
	{
		src: `select
                xxx
                , xxx
                , xxx,
                xxx

        from xxx a
        join xxx b
        on xxx = xxx
        join xxx xxx
        on xxx = xxx
        left outer join xxx
        on xxx = xxx
        where xxx = xxx
        and xxx = true
        and xxx is null
        `,
		want: `
SELECT
  xxx
  , xxx
  , xxx
  , xxx
FROM xxx a
JOIN xxx b
ON xxx = xxx
JOIN xxx xxx
ON xxx = xxx
LEFT OUTER JOIN xxx
ON xxx = xxx
WHERE xxx = xxx
AND xxx = true
AND xxx IS NULL`,
	},
	{
		src: `select
	                count(xxx)
	        from xxx as xxx
	        join xxx on  xxx = xxx
	        where xxx = $1
	        and xxx = $2`,
		want: `
SELECT
  COUNT(xxx)
FROM xxx AS xxx
JOIN xxx
ON xxx = xxx
WHERE xxx = $1
AND xxx = $2`,
	}, {
		src: `select
	                xxx
	        from xxx
	        and xxx in ($3, $4, $5, $6)`,
		want: `
SELECT
  xxx
FROM xxx
AND xxx IN ($3, $4, $5, $6)`,
	},
	{
		src: `select
	                xxx
	                , xxx
	                , xxx
	                , case
	                        when xxx is null THEN $1
	                        ELSE $2
	                  end as xxx
	                , case
	                        when xxx is null THEN 0
	                        ELSE xxx
	                  end as xxx
	        from xxx`,
		want: `
SELECT
  xxx
  , xxx
  , xxx
  , CASE
      WHEN xxx IS NULL THEN $1
      ELSE $2
    END AS xxx
  , CASE
      WHEN xxx IS NULL THEN 0
      ELSE xxx
    END AS xxx
FROM xxx`,
	},
	{
		src: `select
			                xxx
			        from xxx
			        order by xxx`,
		want: `
SELECT
  xxx
FROM xxx
ORDER BY
  xxx`,
	},
	{
		src: `select
	        xxx
	    from xxx
	    where xxx in ($1, $2)
	    order by id`,
		want: `
SELECT
  xxx
FROM xxx
WHERE xxx IN ($1, $2)
ORDER BY
  id`,
	},
	{
		src: `UPDATE xxx
	        SET
	                first_name = $1
	                , last_name = $2
	        WHERE id = $8`,
		want: `
UPDATE
  xxx
SET
  first_name = $1
  , last_name = $2
WHERE id = $8`,
	},
	{
		src: `UPDATE xxx
	        SET
	                xxx = $1
	                , xxx = $2
	        WHERE xxx = $3
	        RETURNING
	          xxx
	          , xxx`,
		want: `
UPDATE
  xxx
SET
  xxx = $1
  , xxx = $2
WHERE xxx = $3
RETURNING
  xxx
  , xxx`,
	},
	{
		src: `SELECT
	                xxx
	        from xxx
	        where xxx = $1
	        and xxx != $3`,
		want: `
SELECT
  xxx
FROM xxx
WHERE xxx = $1
AND xxx != $3`,
	},
	{
		src: `SELECT
	            xxx
	                , xxx
	                , CASE
	                        WHEN xxx IS NULL and xxx IS NULL THEN $1
	                        WHEN xxx IS NULL and xxx IS NOT NULL THEN xxx
	                        ELSE $4
	              END
	                , CASE
	                        WHEN xxx IS NULL xxx IS NULL THEN xxx
	                        WHEN xxx IS NULL THEN xxx
	                        ELSE xxx
	              END AS xxx
	        FROM xxx
	        LEFT OUTER JOIN xxx
	        ON xxx = xxx
	        LEFT OUTER JOIN (
	          select
	               xxx, xxx, xxx  ,xxx
	          from xxx
	          join xxx
	          on xxx = xxx
	          join xxx
	          on xxx = xxx
	          where xxx = $5
	        ) as xxx
	        ON xxx = xxx
	        where xxx`,
		want: `
SELECT
  xxx
  , xxx
  , CASE
      WHEN xxx IS NULL AND xxx IS NULL THEN $1
      WHEN xxx IS NULL AND xxx IS NOT NULL THEN xxx
      ELSE $4
    END
  , CASE
      WHEN xxx IS NULL xxx IS NULL THEN xxx
      WHEN xxx IS NULL THEN xxx
      ELSE xxx
    END AS xxx
FROM xxx
LEFT OUTER JOIN xxx
ON xxx = xxx
LEFT OUTER JOIN (
  SELECT
    xxx
    , xxx
    , xxx
    , xxx
  FROM xxx
  JOIN xxx
  ON xxx = xxx
  JOIN xxx
  ON xxx = xxx
  WHERE xxx = $5
) AS xxx
ON xxx = xxx
WHERE xxx`,
	},
	{
		src: `SELECT
	                xxx
	                , xxx
	                , CASE
	                    WHEN xxx IS NULL AND xxx IS NULL THEN xxx
	                        WHEN xxx > xxx THEN xxx
	                        WHEN xxx <= xxx THEN xxx
	                  END as xxx
	                , xxx
	                , xxx
	                , xxx
	                , xxx
	                , xxx
	                , sum(xxx, 0) as xxx
	                , SUM(xxx, 0) as xxx
	                ,  avg(xxx, 0) as xxx
	                , AVG(xxx, 0) as xxx
	                , max(xxx, 0) as xxx
	                , min(xxx, 0) as xxx
	                , extract(xxx, 0) as xxx
	                , EXTRACT(xxx, 0) as xxx
	                , cast(xxx, 0) as xxx
	                , TRIM(xxx, 0) as xxx
	                , xmlforest(xxx, 0) as xxx
	        FROM xxx
	        LEFT OUTER JOIN (
	                SELECT
	                    xxx
	                        , xxx
	                        , xxx
	                FROM (
	                  select xxx from xxx
	                )
	                WHERE xxx = $1
	        ) xxx
	        ON xxx
	        WHERE xxx = xxx`,
		want: `
SELECT
  xxx
  , xxx
  , CASE
      WHEN xxx IS NULL AND xxx IS NULL THEN xxx
      WHEN xxx > xxx THEN xxx
      WHEN xxx <= xxx THEN xxx
    END AS xxx
  , xxx
  , xxx
  , xxx
  , xxx
  , xxx
  , SUM(xxx, 0) AS xxx
  , SUM(xxx, 0) AS xxx
  , AVG(xxx, 0) AS xxx
  , AVG(xxx, 0) AS xxx
  , MAX(xxx, 0) AS xxx
  , MIN(xxx, 0) AS xxx
  , EXTRACT(xxx, 0) AS xxx
  , EXTRACT(xxx, 0) AS xxx
  , CAST(xxx, 0) AS xxx
  , TRIM(xxx, 0) AS xxx
  , XMLFOREST(xxx, 0) AS xxx
FROM xxx
LEFT OUTER JOIN (
  SELECT
    xxx
    , xxx
    , xxx
  FROM (
    SELECT
      xxx
    FROM xxx
  )
  WHERE xxx = $1
) xxx
ON xxx
WHERE xxx = xxx`,
	},
	{
		src: `select 1 + 1, 2 - 1, 3 * 2, 8 / 2,
	    1 + 1 * 3, 3 + 8 / 7,
	    1+1*3, 312+8/7,
	    4%3, 7^5 from xxx`,
		want: `
SELECT
  1 + 1
  , 2 - 1
  , 3 * 2
  , 8 / 2
  , 1 + 1 * 3
  , 3 + 8 / 7
  , 1+1*3
  , 312+8/7
  , 4%3
  , 7^5
FROM xxx`,
	},
	{
		src: `select
		    array[],
		    array[1]
		  from
		    baz`,
		want: `
SELECT
  array []
  , array [1]
FROM baz`,
	},

	{
		src: `select
		    foo,
		    array ( select
		      bar
		    from
		      quz
		    where
		      baz.foo = quz.foo
		    )
		  from
		    baz`,
		want: `
SELECT
  foo
  , array (
    SELECT
      bar
    FROM quz
    WHERE baz.foo = quz.foo
  )
FROM baz`,
	},
	{
		src: `
	    select
	  '{1,2,3}'::int[],
	  '{{1,2}, {3,4}}'::int[][],
	  '{{1,2}, {3,4}}'::int[][2]
	       from xxx
	    `,
		want: `
SELECT
  '{1,2,3}'::int []
  , '{{1,2}, {3,4}}'::int [] []
  , '{{1,2}, {3,4}}'::int [] [2]
FROM xxx`,
	},
	{
		src: `
	   select
	  '2015-01-01 00:00:00-09'::timestamptz at time zone 'America/Chicago'
	  from xxx`,
		want: `
SELECT
  '2015-01-01 00:00:00-09'::timestamptz AT TIME ZONE 'America/Chicago'
FROM xxx`,
	},
	{
		src: `select
	    foo between bexpr::text and bar,
	    foo between -42 and bar,
	    foo between +3 and bar,
	    foo between 1 + 1 and bar,
	    foo between 1 - 1 and bar,
	    foo between 1 * 1 and bar,
	    foo between 1 / 1 and bar,
	    foo between 1 % 1 and bar,
	    foo between 1 ^ 1 and bar,
	    foo between 1 < 1 and bar,
	    foo between 1 > 1 and bar,
	    foo between 1 = 1 and bar,
	    foo between 1 <= 1 and bar,
	    foo between 1 >= 1 and bar,
	    foo between 1 != 1 and bar,
	    foo between 1 @> 1 and bar,
	    foo between @1 and bar,
	    foo between 5 ! and bar,
	    false between foo is document and bar,
	    false between foo is not document and bar
	  from
	    baz`,
		want: `
SELECT
  foo BETWEEN bexpr::text AND bar
  , foo BETWEEN -42 AND bar
  , foo BETWEEN +3 AND bar
  , foo BETWEEN 1 + 1 AND bar
  , foo BETWEEN 1 - 1 AND bar
  , foo BETWEEN 1 * 1 AND bar
  , foo BETWEEN 1 / 1 AND bar
  , foo BETWEEN 1 % 1 AND bar
  , foo BETWEEN 1 ^ 1 AND bar
  , foo BETWEEN 1 < 1 AND bar
  , foo BETWEEN 1 > 1 AND bar
  , foo BETWEEN 1 = 1 AND bar
  , foo BETWEEN 1 <= 1 AND bar
  , foo BETWEEN 1 >= 1 AND bar
  , foo BETWEEN 1 != 1 AND bar
  , foo BETWEEN 1 @> 1 AND bar
  , foo BETWEEN @1 AND bar
  , foo BETWEEN 5 ! AND bar
  , false BETWEEN foo IS document AND bar
  , false BETWEEN foo IS NOT document AND bar
FROM baz`,
	},
	{
		src: `
	    select
	  b'10101',
	  x'0123456789abcdefABCDEF' from xxx
	    `,
		want: `
SELECT
  b '10101'
  , x '0123456789abcdefABCDEF'
FROM xxx`,
	},
	{
		src: `
	    select
	  foo and bar,
	  baz or quz
	from
	  t
	  `,
		want: `
SELECT
  foo AND bar
  , baz OR quz
FROM t`,
	},
	{
		src: `
	    select
	  not foo,
	  not true,
	  not false,
	  case
	  when foo = bar then
	    7
	  when foo > bar then
	    42
	  else
	    1
	  end
	from
	  t
	    `,
		want: `
SELECT
  NOT foo
  , NOT true
  , NOT false
  , CASE
      WHEN foo = bar THEN 7
      WHEN foo > bar THEN 42
      ELSE 1
    END
FROM t`,
	},
	{
		src: `
	    select
	  case foo
	  when 4 then
	    'A'
	  when 3 then
	    'B'
	  else
	    'C'
	  end
	from
	  baz
	    `,
		want: `
SELECT
  CASE foo
     WHEN 4 THEN 'A'
     WHEN 3 THEN 'B'
     ELSE 'C'
  END
FROM baz`,
	},
	{
		src: `
	    select
	    CAST('{1,2,3}' as int[]),
	    'Foo' collate "C",
	    'Bar' collate "en_US"
	    from xxx`,
		want: `
SELECT
  CAST('{1,2,3}' AS int [])
  , 'Foo' COLLATE "C"
  , 'Bar' COLLATE "en_US"
FROM xxx`,
	},
	{
		src: `select
	    1 = 1,
	    2 > 1,
	    2 < 8,
	    1 != 2,
	    1 != 2,
	    3 >= 2,
	    2 <= 7
	    from xxx
	  `,
		want: `
SELECT
  1 = 1
  , 2 > 1
  , 2 < 8
  , 1 != 2
  , 1 != 2
  , 3 >= 2
  , 2 <= 7
FROM xxx`,
	},
	{
		src: `
SELECT
  CHAR 'hi'
  , CHAR(2) 'hi'
  , VARCHAR 'hi'
  , VARCHAR(2) 'hi'
  , TIMESTAMP(4) '2000-01-01 00:00:00'`,
		want: `
SELECT
  CHAR 'hi'
  , CHAR(2) 'hi'
  , VARCHAR 'hi'
  , VARCHAR(2) 'hi'
  , TIMESTAMP(4) '2000-01-01 00:00:00'`,
	},
	{
		src: `
	    select
	  foo @> bar,
	  @foo,
	  'foo' || 'bar'
	    `,
		want: `
SELECT
  foo @> bar
  , @foo
  , 'foo' || 'bar'`,
	},
	{
		src: `
	    select distinct
	  foo,
	  bar
	from
	  baz
	    `,
		want: `
SELECT DISTINCT
  foo
  , bar
FROM baz`,
	},
	{
		src: `select
	    foo,
	    bar
	  from
	    baz
	  except
	  select
	    a,
	    b
	  from
	    quz
	    `,
		want: `
SELECT
  foo
  , bar
FROM baz
EXCEPT
SELECT
  a
  , b
FROM quz`,
	},
	{
		src: `select
	    foo,
	    bar
	  from
	    baz
	  where
	    exists  ( select
	      1
	    from
	      quz
	    )
	  `,
		want: `
SELECT
  foo
  , bar
FROM baz
WHERE EXISTS (
  SELECT
    1
  FROM quz
)`,
	},
	{
		src: `select
	    extract(year from '2000-01-01 12:34:56'::timestamptz),
	    extract(month from '2000-01-01 12:34:56'::timestamptz)
	  `,
		want: `
SELECT
  EXTRACT(year FROM '2000-01-01 12:34:56'::timestamptz)
  , EXTRACT(month FROM '2000-01-01 12:34:56'::timestamptz)`,
	},
	{
		src: `select
	    coalesce(a, b, c),
	    greatest(d, e, f),
	    least(g, h, i),
	    xmlconcat(j, k, l)
	  from
	    foo
	  `,
		want: `
SELECT
  COALESCE(a, b, c)
  , GREATEST(d, e, f)
  , LEAST(g, h, i)
  , XMLCONCAT(j, k, l)
FROM foo`,
	},
	{
		src: `select
	    foo,
	    bar
	  from
	    baz
	  intersect
	  select
	    a,
	    b
	  from
	    quz
	  `,
		want: `
SELECT
  foo
  , bar
FROM baz
INTERSECT
SELECT
  a
  , b
FROM quz`,
	},
	{
		src: `select
	    interval '5',
	    interval '5' hour,
	    interval '5' hour to minute,
	    interval '5' second(5),
	    interval(2) '10.324'
	  `,
		want: `
SELECT
  INTERVAL '5'
  , INTERVAL '5' hour
  , INTERVAL '5' hour to minute
  , INTERVAL '5' SECOND(5)
  , INTERVAL(2) '10.324'`,
	},
	{
		src: `select
	    foo,
	    bar
	  from
	    baz
	  where
	    foo like 'abd%'
	    or foo like 'ada%' escape '!'
	  `,
		want: `
SELECT
  foo
  , bar
FROM baz
WHERE foo LIKE 'abd%'
OR foo LIKE 'ada%' escape '!'`,
	},
	{
		src: `select
	  foo,
	  bar
	from
	  baz
	limit 7
	offset 42`,
		want: `
SELECT
  foo
  , bar
FROM baz
LIMIT 7
OFFSET 42`,
	},
	{
		src: `select foo, bar from baz offset 42 rows fetch next 7 rows only
	    `,
		want: `
SELECT
  foo
  , bar
FROM baz
OFFSET 42 ROWS
FETCH next 7 ROWS only`,
	},
	{
		src: `select
	    foo,
	    bar
	  from
	    baz
	  order by
	    foo desc nulls first,
	    quz asc nulls last,
	    abc nulls last
	  `,
		want: `
SELECT
  foo
  , bar
FROM baz
ORDER BY
  foo DESC NULLS FIRST
  , quz ASC NULLS LAST
  , abc NULLS LAST`,
	},
	{
		src: `select
	    (date '2000-01-01', date '2000-01-31') overlaps (date '2000-01-15', date '2000-02-15')
	  `,
		want: `
SELECT
  (date '2000-01-01', date '2000-01-31') OVERLAPS (date '2000-01-15', date '2000-02-15')`,
	},
	{
		src: `select
	    foo,
	    row_number() over(range unbounded preceding)
	  from
	    baz
	  `,
		want: `
SELECT
  foo
  , ROW_NUMBER() OVER(range unbounded preceding)
FROM baz`,
	},
	{
		src: `select xxx from xxx union all select xxx from xxx`,
		want: `
SELECT
  xxx
FROM xxx
UNION ALL
SELECT
  xxx
FROM xxx`,
	},
	{
		src: `lock table in xxx`,
		want: `
LOCK table
IN xxx`,
	},
}
