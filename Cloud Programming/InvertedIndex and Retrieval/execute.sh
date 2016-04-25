hdfs dfs -rm -r /user/s101062115/InvertedIndex/output
hdfs dfs -rm -r /user/s101062115/Query/output

hadoop jar InvertedIndex.jar invertedIndex.InvertedIndex /user/s101062115/input /user/s101062115/InvertedIndex/output
hadoop jar Query.jar query.Query /user/s101062115/InvertedIndex/output /user/s101062115/input /user/s101062115/Query/output
hdfs dfs -cat /user/s101062115/Query/output/part-*
