package invertedIndex;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class InvertedIndexKeyComparator extends WritableComparator {

	public InvertedIndexKeyComparator() {
		super(InvertedIndexKeyPair.class, true);
	}


	public int compare(WritableComparable o1, WritableComparable o2) {
		InvertedIndexKeyPair key1 = (InvertedIndexKeyPair) o1;
		InvertedIndexKeyPair key2 = (InvertedIndexKeyPair) o2;

		int result =  key1.getTerm().compareTo(key2.getTerm());
		if(result == 0) {
			result = key1.getDocID() - key2.getDocID();
		}

		return result;
	}
}
