package query;

import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class QueryKeyComparator extends WritableComparator {

	public QueryKeyComparator() {
		super(QueryKeyPair.class, true);
	}

	public int compare(WritableComparable o1, WritableComparable o2) {
		QueryKeyPair key1 = (QueryKeyPair) o1;
		QueryKeyPair key2 = (QueryKeyPair) o2;

		int result = key1.getWeight().compareTo(key2.getWeight());
		if(result == 0) {
			result = key1.getDocID() - key2.getDocID();
		}
		else result *= -1;
		return result;
	}
}
