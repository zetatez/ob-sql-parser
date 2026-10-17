package token

// Keyword classification.
type Keyword struct {
	Reserved bool
}

// keywords maps the normalized (lowercase) word to its classification.
//
// This is an intentionally small but representative table. It is meant to grow
// toward the full set of OB 4.x MySQL-mode reserved/non-reserved keywords as
// statement coverage expands. `Reserved:true` words cannot be used as bare
// identifiers; `Reserved:false` words are non-reserved and may.
var keywords = map[string]Keyword{

	"select":        {Reserved: true},
	"insert":        {Reserved: true},
	"update":        {Reserved: true},
	"delete":        {Reserved: true},
	"create":        {Reserved: true},
	"alter":         {Reserved: true},
	"drop":          {Reserved: true},
	"truncate":      {Reserved: false},
	"from":          {Reserved: true},
	"where":         {Reserved: true},
	"group":         {Reserved: true},
	"by":            {Reserved: true},
	"order":         {Reserved: true},
	"having":        {Reserved: true},
	"limit":         {Reserved: false},
	"union":         {Reserved: true},
	"except":        {Reserved: true},
	"intersect":     {Reserved: true},
	"minus":         {Reserved: true},
	"with":          {Reserved: true},
	"recursive":     {Reserved: false},
	"window":        {Reserved: false},
	"inner":         {Reserved: true},
	"outer":         {Reserved: true},
	"left":          {Reserved: true},
	"right":         {Reserved: true},
	"full":          {Reserved: true},
	"cross":         {Reserved: true},
	"join":          {Reserved: true},
	"on":            {Reserved: true},
	"using":         {Reserved: true},
	"natural":       {Reserved: true},
	"straight_join": {Reserved: true},
	"as":            {Reserved: true},
	"distinct":      {Reserved: true},
	"all":           {Reserved: true},
	"and":           {Reserved: true},
	"or":            {Reserved: true},
	"xor":           {Reserved: true},
	"not":           {Reserved: true},
	"is":            {Reserved: true},
	"null":          {Reserved: false},
	"like":          {Reserved: true},
	"regexp":        {Reserved: true},
	"rlike":         {Reserved: true},
	"in":            {Reserved: true},
	"between":       {Reserved: true},
	"case":          {Reserved: true},
	"when":          {Reserved: true},
	"then":          {Reserved: true},
	"else":          {Reserved: true},
	"end":           {Reserved: true},
	"exists":        {Reserved: false},
	"binary":        {Reserved: false},
	"cast":          {Reserved: false},
	"interval":      {Reserved: false},
	"escaped":       {Reserved: false},
	"table":         {Reserved: true},
	"database":      {Reserved: false},
	"schema":        {Reserved: false},
	"index":         {Reserved: false},
	"key":           {Reserved: false},
	"primary":       {Reserved: true},
	"foreign":       {Reserved: true},
	"references":    {Reserved: true},
	"constraint":    {Reserved: false},
	"unique":        {Reserved: false},
	"default":       {Reserved: true},
	"into":          {Reserved: true},
	"values":        {Reserved: true},
	"value":         {Reserved: true},
	"set":           {Reserved: true},
	"show":          {Reserved: false},
	"use":           {Reserved: false},
	"explain":       {Reserved: false},
	"begin":         {Reserved: false},
	"commit":        {Reserved: false},
	"rollback":      {Reserved: false},
	"desc":          {Reserved: false},
	"asc":           {Reserved: false},
	"nulls":         {Reserved: false},
	"first":         {Reserved: false},
	"last":          {Reserved: false},
	"over":          {Reserved: false},
	"partition":     {Reserved: false},
	"ignore":        {Reserved: false},
	"replace":       {Reserved: true},
	"duplicate":     {Reserved: false},
	"low_priority":  {Reserved: false},
	"high_priority": {Reserved: false},
	"delayed":       {Reserved: false},
	"quick":         {Reserved: false},
	"optimize":      {Reserved: false},
	"repair":        {Reserved: false},
	"analyze":       {Reserved: false},
	"check":         {Reserved: false},
	"flush":         {Reserved: false},
	"rename":        {Reserved: false},
	"alter_system":  {Reserved: false},
	"zero":          {Reserved: false},
	"system":        {Reserved: false},
	"epoch":         {Reserved: false},

	"character":         {Reserved: false},
	"charset":           {Reserved: false},
	"comment":           {Reserved: false},
	"collate":           {Reserved: false},
	"current_timestamp": {Reserved: false},
	"engine":            {Reserved: false},
	"fulltext":          {Reserved: false},
	"generated":         {Reserved: false},
	"unsigned":          {Reserved: false},
	"timestamp":         {Reserved: false},
	"blob":              {Reserved: false},
	"text":              {Reserved: false},
	"int":               {Reserved: false},
	"decimal":           {Reserved: false},
	"double":            {Reserved: false},
	"bigint":            {Reserved: false},
	"smallint":          {Reserved: false},
	"tinyint":           {Reserved: false},
	"float":             {Reserved: false},
	"enum":              {Reserved: false},
	"json":              {Reserved: false},
	"year":              {Reserved: false},
	"date":              {Reserved: false},
	"time":              {Reserved: false},
	"datetime":          {Reserved: false},
	"varchar":           {Reserved: false},
	"char":              {Reserved: false},
	"boolean":           {Reserved: false},
	"number":            {Reserved: false},
	"row":               {Reserved: false},
	"rows":              {Reserved: false},
	"range":             {Reserved: false},
	"hash":              {Reserved: false},
	"list":              {Reserved: false},
	"columns":           {Reserved: false},
	"column":            {Reserved: false},
	"separator":         {Reserved: false},
	"views":             {Reserved: false},
	"status":            {Reserved: false},
	"variables":         {Reserved: false},
	"global":            {Reserved: false},
	"session":           {Reserved: false},
	"local":             {Reserved: false},
	"tables":            {Reserved: false},
	"databases":         {Reserved: false},
	"schemas":           {Reserved: false},
	"warnings":          {Reserved: false},
	"errors":            {Reserved: false},
	"create_statement":  {Reserved: false},
	"processlist":       {Reserved: false},
	"profiles":          {Reserved: false},
	"triggers":          {Reserved: false},
	"procedure":         {Reserved: false},
	"function":          {Reserved: false},
	"grants":            {Reserved: false},
	"collation":         {Reserved: false},
	"open":              {Reserved: false},
	"major":             {Reserved: false},
	"minor":             {Reserved: false},
	"read":              {Reserved: false},
	"write":             {Reserved: false},
	"dual":              {Reserved: false},
	"current_user":      {Reserved: false},
	"user":              {Reserved: false},
	"password":          {Reserved: false},

	"current":     {Reserved: false},
	"unbounded":   {Reserved: false},
	"preceding":   {Reserved: false},
	"following":   {Reserved: false},
	"names":       {Reserved: false},
	"work":        {Reserved: false},
	"transaction": {Reserved: false},
	"only":        {Reserved: false},
	"chain":       {Reserved: false},
	"temporary":   {Reserved: false},
	"external":    {Reserved: false},
	"restrict":    {Reserved: false},
	"cascade":     {Reserved: false},
	"extended":    {Reserved: false},
	"pretty":      {Reserved: false},
	"format":      {Reserved: false},
	"share":       {Reserved: false},
	"sounds":      {Reserved: false},
	"div":         {Reserved: false},
	"mod":         {Reserved: false},
	"fields":      {Reserved: false},
	"engines":     {Reserved: false},
	"storage":     {Reserved: false},
	"indexes":     {Reserved: false},
	"keys":        {Reserved: false},
}

// LookupKeyword returns whether the word is a known keyword and whether it is
// reserved. word is matched case-insensitively (normalized to lowercase).
func LookupKeyword(word string) (isKeyword bool, reserved bool) {
	k, ok := keywords[toLower(word)]
	if !ok {
		return false, false
	}
	return true, k.Reserved
}

// IsReservedWord reports whether the word is a reserved keyword.
func IsReservedWord(word string) bool {
	_, r := LookupKeyword(word)
	return r
}

// IsKeyword reports whether the word is a known keyword (reserved or not).
func IsKeyword(word string) bool {
	ok, _ := LookupKeyword(word)
	return ok
}
func toLower(b string) string {
	out := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		c := b[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}
