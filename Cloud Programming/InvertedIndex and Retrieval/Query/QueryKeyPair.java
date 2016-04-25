package query;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

import org.apache.hadoop.io.DoubleWritable;
import org.apache.hadoop.io.WritableComparable;

public class QueryKeyPair implements WritableComparable<QueryKeyPair> {

    private DoubleWritable weight;
    private int docID;

    public QueryKeyPair() {
	this.weight = new DoubleWritable();
    }

    public QueryKeyPair(DoubleWritable weight, int docID) {
	//TODO: constructor
	this.weight = weight;
	this.docID = docID;
    }

    @Override
	public void write(DataOutput out) throws IOException {
	    weight.write(out);
	    out.writeInt(docID);
	}

    @Override
	public void readFields(DataInput in) throws IOException {
	    weight.readFields(in);
	    docID = in.readInt();
	}
    @Override
	public int compareTo(QueryKeyPair o) {
	    return weight.compareTo(o.weight);
	}

    public String getWeight() {
	return this.weight.toString();
    }

    public int getDocID() {
	return this.docID;
    }

}
