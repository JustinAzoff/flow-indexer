package loggen

import (
	"bytes"
	"fmt"
	"math/rand"
)

func RandomIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func PartiallyRandomIPv4(randomOctects int) string {
	switch randomOctects {
	case 0:
		return "1.2.3.4"
	case 1:
		return fmt.Sprintf("%d.2.3.4", rand.Intn(256))
	case 2:
		return fmt.Sprintf("%d.%d.3.4", rand.Intn(256), rand.Intn(256))
	case 3:
		return fmt.Sprintf("%d.%d.%d.4", rand.Intn(256), rand.Intn(256), rand.Intn(256))
	case 4:
		return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
	}
	return "1.2.3.4"
}

var template = "1324071333.493287	ChcmbpkpM9NGNEFVi	%s	51880	%s	22	tcp	ssh	6.159326	2669	2501	SF	T	0	ShAdDaFf	25	3981	20	3549	(empty)\n"

func RandomASCIIBroLog(lines int) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("#some random lines\n#with comments\n#doesn't matter\n")
	for i := 0; i < lines; i++ {
		fmt.Fprintf(buf, template, RandomIPv4(), RandomIPv4())
	}
	return buf.Bytes()
}

var jsTemplate = `{"ts":1409516196.337184,"uid":"CBQ0o7d3H96U3UvZi","id.orig_h":"%s","id.orig_p":40184,"id.resp_h":"%s","id.resp_p":41644,"proto":"tcp","service":"ssh","duration":0.392307,"orig_bytes":3205,"resp_bytes":2129,"conn_state":"S1","missed_bytes":0,"history":"ShADad","orig_pkts":12,"orig_ip_bytes":3837,"resp_pkts":12,"resp_ip_bytes":2761,"tunnel_parents":[]}`

func RandomJSONBroLog(lines int) []byte {
	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < lines; i++ {
		fmt.Fprintf(buf, jsTemplate, RandomIPv4(), RandomIPv4())
		fmt.Fprint(buf, "\n")
	}
	return buf.Bytes()
}
