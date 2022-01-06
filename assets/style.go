// Code generated by go generate; DO NOT EDIT.

package assets

var AssetsMap = map[string]string{
	"style": `body {
	max-width: 750px;
	margin: auto;
	padding: 5px;
}

/* Lists ************************************************************/

ol.posts, ol.replies {
	padding: 0;
	list-style: none;
}

ol.posts > li:not(:last-child),
ol.replies > li:not(:last-child) {
	margin-bottom: 1em;
}

/* Posts ************************************************************/

table.posts {
	width: 100%;
	border-collapse: collapse;
	border: 1px solid darkgrey;
}
th {
	background-color: darkgrey;
}
th, td {
	padding: .2em;
}
.posts-title {
	width: 100%;
}
td {
	background-color: #eee;
	border: 1px solid darkgrey;
}
.posts h2 {
	margin: 0;
	font-size: 1em;
}

.content h1 { font-size: 1.5em; }
.content h2 { font-size: 1.2em; }
.content h3 { font-size: 1em; }
.content { margin: 1em 0; }

/* Replies **********************************************************/

.reply {
	border-left: 1px solid lightgrey;
	display: block;
}

.reply summary {
	background-color: lightgrey;
	display: list-item;
	padding: .2em;
}

.reply .content, .reply footer {
	padding: 0 1em;
}

.thread {
	padding: 1em 0 1em 1em;
}

.meta {
	background-color: lightgrey;
	padding: .2em;
}

/* Topics ***********************************************************/

.topics {
	background-color: lightgrey;
	padding: .2em;
}

.topics .selected {
	font-weight: bold;
}

/* Navigation *******************************************************/

header > nav {
	float: right;
}

body > footer {
	margin-top: 1em;
	border-top: 1px solid lightgrey;
	color: grey;
	text-align: center;
}

/* Forms ************************************************************/

.auth-form {
	max-width: 200px;
}

.field {
	margin-bottom: 1em;
}

.field label {
	display: block;
}

input[type=text], input[type=password] {
	width: 100%;
	box-sizing: border-box;
}

textarea {
	width: 100%;
	height: 250px;
	display: block;
	box-sizing: border-box;
}

/* Misc *************************************************************/

.key-value {
	padding: 0;
	margin: 0;
	list-style: none;
}

blockquote {
	margin: 0;
	color: green;
	font-style: italic;
}`,
}
