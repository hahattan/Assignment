package invertedIndex;

import java.io.IOException;

import org.apache.hadoop.io.DoubleWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

public class InvertedIndexReducer extends Reducer<InvertedIndexKeyPair,InvertedIndexValuePair,Text,Text> {

  private Text result = new Text();

  public void reduce(InvertedIndexKeyPair key, Iterable<InvertedIndexValuePair> values, Context context) throws IOException, InterruptedException {
    int count = 0;
    StringBuilder builder = new StringBuilder();
    // agrregate the amount of same starting character
    for (InvertedIndexValuePair val: values) {
      // count += ...
      count += 1;
      builder.append(val.getDocID() + " ");
      builder.append(val.getTermFrequency() + " ");
      builder.append(val.getOffset() + ";");
    }
    StringBuilder table = new StringBuilder();
    table.append(String.valueOf(count));
    table.append(";");
    table.append(builder.toString());
    result.set(table.toString());
    // write the result
    context.write(new Text(key.getTerm()), result);

  }
}
