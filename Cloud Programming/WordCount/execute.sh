# Do not uncomment these lines to directly execute the script
# Modify the path to fit your need before using this script
#hdfs dfs -rm -r /user/TA/WordCount/Output/
#hadoop jar WordCount.jar wordcount.WordCount /user/shared/WordCount/Input /user/TA/WordCount/Output
#hdfs dfs -cat /user/TA/WordCount/Output/part-*

hdfs dfs -rm -r /user/s101062115/WordCount/output/
hadoop jar WordCount.jar wordcount.WordCount /user/s101062115/input1 /user/s101062115/WordCount/output
hdfs dfs -cat /user/s101062115/WordCount/output/part-*
