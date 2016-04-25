package invertedIndex;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Partitioner;

public class InvertedIndexPartitioner extends Partitioner<InvertedIndexKeyPair, InvertedIndexValuePair> {
  @Override
    public int getPartition(InvertedIndexKeyPair key, InvertedIndexValuePair values, int numReduceTasks) {
      // customize which <K ,V> will go to which reducer
      /*
      if((key.getTerm().charAt(0) >= 65 && key.getTerm().charAt(0) <= 71) || ((key.getTerm().charAt(0) >= 97 && key.getTerm().charAt(0) <= 103))) {
	return 0;
      }else {
	return 1;
      }
      */
      
      
      int hash = key.getTerm().hashCode();
      int partition = (hash & Integer.MAX_VALUE) % numReduceTasks;
      return partition;
      
    }
}
