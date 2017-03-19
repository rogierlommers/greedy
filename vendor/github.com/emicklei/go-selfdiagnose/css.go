package selfdiagnose

// css has the contains of the style.css file
func css() string {
	return `
body, table {
	font-family: 'Lucida Sans Typewriter', 'Lucida Console', monaco, 'Bitstream Vera Sans Mono', monospace;
	font-size: 10px;
}
.odd {
	background-color: #E6ECE9;
}
.even {
	background-color: #DBE4DF;
}

.odd.failed.warning {
	background-color: #FFC181;
}
.even.failed.warning {
	background-color: #FFB364;
}

.odd.failed.critical, .odd.error {
	background-color: #FD9E9E;
}
.even.failed.critical, .even.error {
	background-color: #FF8282;
}
.header {
	background-color: #d6d6d6
}
`
}
