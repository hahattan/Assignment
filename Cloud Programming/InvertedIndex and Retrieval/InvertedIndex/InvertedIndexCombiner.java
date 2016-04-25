package invertedIndex;

import java.io.IOException;
import java.util.*;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.mapreduce.Reducer;

public class InvertedIndexCombiner extends Reducer<InvertedIndexKeyPair, InvertedIndexValuePair,InvertedIndexKeyPair, InvertedIndexValuePair> {
  // Combiner implements method in Reducer

  public void reduce(InvertedIndexKeyPair key, Iterable<InvertedIndexValuePair> values, Context context) throws IOException, InterruptedException {
    StringBuilder builder = new StringBuilder();
    int count = 0;
    List<Integer> offset = new ArrayList<Integer>();
    for (InvertedIndexValuePair val: values) {
      //TODO: agrregate the result from mapper
      count += 1;
      offset.add(Integer.valueOf(val.getOffset()));
    }
    Collections.sort(offset);
    builder.append("[");
    for(int i = 0; i < offset.size(); i++) {
        if(i != 0) builder.append(",");
        builder.append(String.valueOf(offset.get(i)));
    }
    builder.append("]");
    InvertedIndexValuePair temp = new InvertedIndexValuePair(key.getDocID(), new IntWritable(count), new Text(builder.toString()));
    //context.write(K,V)
    context.write(key, temp);

  }
}
