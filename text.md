1 Framing
    1.1 Text Messages
    1.2 Binary Messages
2 Pings/Pongs
3 Reserved Bits
4 Opcodes
    4.1 Non-control Opcodes
    4.2 Control Opcodes
5 Fragmentation
6 UTF-8 Handling
    6.1 Valid UTF-8 with zero payload fragments
    6.2 Valid UTF-8 unfragmented, fragmented on code-points and within code-points
    6.3 Invalid UTF-8 differently fragmented
    6.4 Fail-fast on invalid UTF-8
    6.5 Some valid UTF-8 sequences
    6.6 All prefixes of a valid UTF-8 string that contains multi-byte code points
    6.7 First possible sequence of a certain length
    6.8 First possible sequence length 5/6 (invalid codepoints)
    6.9 Last possible sequence of a certain length
    6.10 Last possible sequence length 4/5/6 (invalid codepoints)
    6.11 Other boundary conditions
    6.12 Unexpected continuation bytes
    6.13 Lonely start characters
    6.14 Sequences with last continuation byte missing
    6.15 Concatenation of incomplete sequences
    6.16 Impossible bytes
    6.17 Examples of an overlong ASCII character
    6.18 Maximum overlong sequences
    6.19 Overlong representation of the NUL character
    6.20 Single UTF-16 surrogates
    6.21 Paired UTF-16 surrogates
    6.22 Non-character code points (valid UTF-8)
    6.23 Unicode specials (i.e. replacement char)
7 Close Handling
    7.1 Basic close behavior (fuzzer initiated)
    7.3 Close frame structure: payload length (fuzzer initiated)
    7.5 Close frame structure: payload value (fuzzer initiated)
    7.7 Close frame structure: valid close codes (fuzzer initiated)
    7.9 Close frame structure: invalid close codes (fuzzer initiated)
    7.13 Informational close information (fuzzer initiated)

9 Limits/Performance
    9.1 Text Message (increasing size)
    9.2 Binary Message (increasing size)
    9.3 Fragmented Text Message (fixed size, increasing fragment size)
    9.4 Fragmented Binary Message (fixed size, increasing fragment size)
    9.5 Text Message (fixed size, increasing chop size)
    9.6 Binary Text Message (fixed size, increasing chop size)
    9.7 Text Message Roundtrip Time (fixed number, increasing size)
    9.8 Binary Message Roundtrip Time (fixed number, increasing size)
10 Misc
    10.1 Auto-Fragmentation

12 WebSocket Compression (different payloads)
    12.1 Large JSON data file (utf8, 194056 bytes)
    12.2 Lena Picture, Bitmap 512x512 bw (binary, 263222 bytes)
    12.3 Human readable text, Goethe's Faust I (German) (binary, 222218 bytes)
    12.4 Large HTML file (utf8, 263527 bytes)
    12.5 A larger PDF (binary, 1042328 bytes)

13 WebSocket Compression (different parameters)
    13.1 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(False, 0)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(False, 0)]
    13.2 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 0)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 0)]
    13.3 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(False, 9)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(False, 9)]
    13.4 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(False, 15)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(False, 15)]
    13.5 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 9)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 9)]
    13.6 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 15)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 15)]
    13.7 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 9), (True, 0), (False, 0)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 9), (True, 0), (False, 0)]
