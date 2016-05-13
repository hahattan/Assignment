DROP TABLE IF EXISTS math;
DROP TABLE IF EXISTS eng;
DROP TABLE IF EXISTS s101062115_score;

CREATE EXTERNAL TABLE math
(key string, name string, grade int)
STORED BY 'org.apache.hadoop.hive.hbase.HBaseStorageHandler'
WITH SERDEPROPERTIES ("hbase.columns.mapping" = ":key,grade:name,grade:math")
TBLPROPERTIES ("hbase.table.name" = "s101062115:math");

CREATE EXTERNAL TABLE eng
(key string, name string, grade int)
STORED BY 'org.apache.hadoop.hive.hbase.HBaseStorageHandler'
WITH SERDEPROPERTIES ("hbase.columns.mapping" = ":key,grade:name,grade:eng")
TBLPROPERTIES ("hbase.table.name" = "s101062115:eng");

CREATE TABLE s101062115_score AS
SELECT math.key AS name, math.grade AS math, eng.grade AS eng, (math.grade+eng.grade)/2 AS average
FROM math JOIN eng ON(math.key = eng.key);

SELECT COUNT(*)
FROM s101062115_score
WHERE average < 60;

SELECT *
FROM s101062115_score
ORDER BY average DESC, name ASC
LIMIT 5;
