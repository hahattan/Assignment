package query;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.Writable;

public class QueryValuePair implements Writable {

  private int docID;
  private Text offset;

  public QueryValuePair() {
    this.offset = new Text();
  }

  public QueryValuePair(int docID, Text offset) {
    //TODO: constructor
    this.docID = docID;
    this.offset = offset;
  }

  @Override
    public void write(DataOutput out) throws IOException {
      out.writeInt(docID);
      offset.write(out);
    }

  @Override
    public void readFields(DataInput in) throws IOException {
      docID = in.readInt();
      offset.readFields(in);
    }

    public int getDocID() {
      return this.docID;
    }

  public String getOffset() {
      return this.offset.toString();
  }

}
