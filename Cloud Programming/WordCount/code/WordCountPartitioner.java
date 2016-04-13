package wordcount;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Partitioner;

public class WordCountPartitioner extends Partitioner<Text, IntWritable> {
  @Override
    public int getPartition(Text key, IntWritable value, int numReduceTasks) {
      // customize which <K ,V> will go to which reducer
      if((key.charAt(0)>=65 && key.charAt(0) <= 71) || (key.charAt(0)>=97 && key.charAt(0)<=103)) return 0;
      else return 1;
    }
}
