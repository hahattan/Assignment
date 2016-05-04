hdfs dfs -rm -r /user/s101062115/InvertedIndex
hdfs dfs -rm -r /user/s101062115/Query

hadoop jar InvertedIndex.jar invertedIndex.InvertedIndex input /user/s101062115/InvertedIndex
hadoop jar Query.jar query.Query /user/s101062115/InvertedIndex input /user/s101062115/Query
hdfs dfs -cat /user/s101062115/Query/part-*
