package query;

import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class QueryGroupComparator extends WritableComparator {

    public QueryGroupComparator() {
        super(QueryKeyPair.class, true);
    }


    public int compare(WritableComparable o1, WritableComparable o2) {
        QueryKeyPair key1 = (QueryKeyPair) o1;
        QueryKeyPair key2 = (QueryKeyPair) o2;

        return key1.getWeight().compareTo(key2.getWeight());
    }
}
