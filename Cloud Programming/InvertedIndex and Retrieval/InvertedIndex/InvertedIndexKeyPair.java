package invertedIndex;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.WritableComparable;

public class InvertedIndexKeyPair implements WritableComparable<InvertedIndexKeyPair> {


    private Text term;
    private int docID;

    public InvertedIndexKeyPair() {
        this.term = new Text();
    }

    public InvertedIndexKeyPair(Text term, int docID) {
        //TODO: constructor
        this.term = term;
        this.docID = docID;
    }

    @Override
    public void write(DataOutput out) throws IOException {
        term.write(out);
        out.writeInt(docID);
    }

    @Override
    public void readFields(DataInput in) throws IOException {
        term.readFields(in);
        docID = in.readInt();
    }

    @Override
    public int compareTo(InvertedIndexKeyPair o) {
        return term.compareTo(o.term);
    }

    public String getTerm() {
        return this.term.toString();
    }

    public int getDocID() {
        return this.docID;
    }
}
