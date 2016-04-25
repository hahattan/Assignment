package invertedIndex;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.Writable;
import org.apache.hadoop.io.IntWritable;

public class InvertedIndexValuePair implements Writable {

  private int docID;
  private IntWritable tf;
  private Text offset;

  public InvertedIndexValuePair() {
    this.tf = new IntWritable();
    this.offset = new Text();
  }

  public InvertedIndexValuePair(int docID, IntWritable tf, Text offset) {
    //TODO: constructor
    this.docID = docID;
    this.tf = tf;
    this.offset = offset;
  }

  @Override
    public void write(DataOutput out) throws IOException {
      out.writeInt(docID);
      tf.write(out);
      offset.write(out);
    }

  @Override
    public void readFields(DataInput in) throws IOException {
      docID = in.readInt();
      tf.readFields(in);
      offset.readFields(in);
    }

    public int getDocID() {
      return this.docID;
    }

  public int getTermFrequency() {
    return this.tf.get();
  }

  public String getOffset() {
      return this.offset.toString();
  }

}
