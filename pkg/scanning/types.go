package scanning

/* XXX -
 * This is a bad enum, there is no indication to what V1 or V2 are referring to,
 * a better way to write this would be:
 *
 * type DataVersion int
 * const (
 * 	 DataVersion_V1 DataVersion = iota
 *	 DataVersion_V2
 * )
 *
 * Also, cf. comment below as to why we don't need this.
 */
const (
	Version = iota
	V1
	V2
)

type Scan struct {
	Ip          string      `json:"ip"`

	/* XXX -
	 * Why bother specify a fixed length integer ? There is strictly no
	 * reason to have a 32bits port number, as it is ubiquitously 16bits. If
	 * the intend is to have a generic "unsigned int" type, then 'uint' will
	 * do. If the intend is to have a precise type, then it should be set to
	 * `uint16`, but certainly not `uint32'.
	 */
	Port        uint32      `json:"port"`
	Service     string      `json:"service"`
	Timestamp   int64       `json:"timestamp"`
	DataVersion int         `json:"data_version"`
	Data        interface{} `json:"data"`
}

/* XXX -
 * This form of API does not work well with Golang's encoding/json package, as
 * there is no definite type and the interface{} will end up being decoded in a
 * map[string]interface() incuring a higher runtime cost.
 *
 * Having dedicated fields for the v1 and v2 remove any requirement for a
 * `data_version' entry.
 *
 */
type V1Data struct {
	ResponseBytesUtf8 []byte `json:"response_bytes_utf8"`
}

type V2Data struct {
	ResponseStr string `json:"response_str"`
}
