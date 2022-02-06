// Code generated by go generate; DO NOT EDIT.

package assets

var AssetsMap = map[string]string{
	"style": `body {
	max-width: 800px;
	margin: auto;
	padding: 5px;
	font: 16px/1.4 system-ui, sans-serif;
	background-color: mintcream;
}

/* Lists ************************************************************/
.errors {
	background-color: mistyrose;
	color: red;
}
.info {
	background-color: palegreen;
	color: green;
}
.topic img {
	max-width: 100%;
}

main {
	margin-bottom: 1em;
}
ol.posts, ol.replies {
	padding: 0;
	list-style: none;
}

/*ol.posts > li:not(:last-child),*/
ol.replies > li:not(:last-child) {
	margin-bottom: 1em;
}

/* Posts ************************************************************/



.topic > tbody > tr:nth-child(2n) {
	background-color: whitesmoke;
}
.signature > * {
	margin-bottom: 0;
}
.action {
	margin: 1em 1em 1em 0;
}
nav.breadcrumb {
	padding: 1em;
	border: 1px solid;
	margin: 1em 0;
}
.breadcrumb ul {
	list-style: none;
	padding-left: 1em;
}
.breadcrumb > ul {
	padding: 0;
	margin: 0;
}
.col-author {
	text-align: center;
	width: 100px;
}
table.post {
	background-color: white;
}
table.post .header {
	background-color: paleturquoise;
}

article {
	border: 1px solid;
	background-color: white;
	margin-bottom: 1em;
}

article > header {
	background-color: paleturquoise;
	border-bottom: 1px solid;
	padding: 5px 1em;
}

article > div {
	padding: 0 1em;
}

.sticky, .forum { background-color: cornsilk; }

/* Start post */
/* With a table */
.post-aside {
	text-align: center;
	width: 150px;
	background-color: paleturquoise;
}

.post-body {
	height: 100%;
}

.post-body, .post-body td, .post-body tr {
	padding: 0;
	border: 0;
}

.post-body tbody {
	background-color: inherit;
}
/* With articles */
/*article {*/
/*	padding: 5px;*/
/*	border: 1px solid;*/
/*	margin-bottom: 1em;*/
/*}*/
/*article:after {*/
/*	clear: both;*/
/*	content: "";*/
/*	display: block;*/
/*	visibility: hidden;*/
/*}*/
/*.post-aside {*/
/*	text-align: center;*/
/*	width: 150px;*/
/*	float: left;*/
/*}*/
/*.post-content {*/
/*	max-width: 650px;*/
/*	margin-left: 150px;*/
/*}*/
/* End */
table {
	border-collapse: collapse;
	border: 1px solid;
	width: 100%;
}
tr, td, th {
	vertical-align: top;
	border: 1px solid;
	padding: .5em;
}
thead {
	background-color: paleturquoise;
}
tbody {
	background-color: white;
}


.content h1 { font-size: 1.5em; }
.content h2 { font-size: 1.2em; }
.content h3 { font-size: 1em; }
.content { margin: 1em 0; }

/* Navigation *******************************************************/

header > nav {
	float: right;
}

/*body > footer {*/
/*	margin-top: 1em;*/
/*	border-top: 1px solid lightgrey;*/
/*	color: grey;*/
/*	text-align: center;*/
/*}*/

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

blockquote {
	margin: 0;
	color: green;
	font-style: italic;
}

.center { text-align: center; }
.grow { width: 100%; }

hr {
	border: none;
	height: 1px;
	background-color: grey;
}

.small {
	font-size: 12px;
	color: grey;
}

.col-content {
	display: flex;
	flex-direction: column;
}

.col-content > div {
	flex-grow: 1;
}`,
}
