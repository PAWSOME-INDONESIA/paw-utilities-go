package elastic

import (
	"strings"
)

type (
	Query interface {
		Source() map[string]interface{}
	}
	Queries []Query
)

func (q Queries) Source() []map[string]interface{} {
	var response = make([]map[string]interface{}, 0)
	for _, query := range q {
		response = append(response, query.Source())
	}
	return response
}

type EsQuery struct {
	from   int
	size   int
	source Query
	query  Query
	sort   Queries
}

func NewEsQuery(from, size int, source, query Query, sort Queries) Query {
	return &EsQuery{
		from:   from,
		size:   size,
		source: source,
		query:  query,
		sort:   sort,
	}
}

func (q EsQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	if q.size == 0 {
		response["size"] = 10
	} else {
		response["size"] = q.size
	}
	response["from"] = q.from
	response["_source"] = q.source.Source()
	response["query"] = q.query.Source()
	response["sort"] = q.sort.Source()
	return response
}

type SourceQuery struct {
	includes []string
	excludes []string
}

func NewSourceQuery(includes, excludes []string) Query {
	return &SourceQuery{
		includes: includes,
		excludes: excludes,
	}
}

func (s SourceQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	if len(s.includes) > 0 {
		response["includes"] = s.includes
	}
	if len(s.excludes) > 0 {
		response["excludes"] = s.excludes
	}
	return response
}

type SortQuery struct {
	field string
	order string
}

func NewSortQuery(field, order string) Query {
	return &SortQuery{
		field: field,
		order: order,
	}
}

func (s SortQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	s.order = strings.ToLower(s.order)
	if len(s.order) == 0 || s.order != "desc" {
		s.order = "asc"
	}

	if len(s.field) > 0 {
		response[s.field] = map[string]string{"order": s.order}
	}
	return response
}

type ScriptQuery struct {
	_type  string
	script Query
	order  string
}

func NewScriptQuery(_type, order string, script Query) Query {
	return &ScriptQuery{
		_type:  _type,
		order:  order,
		script: script,
	}
}

func (s ScriptQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	response["_script"] = map[string]interface{}{
		"type":   s._type,
		"script": s.script.Source(),
		"order":  s.order,
	}
	return response
}

type Script struct {
	lang         string
	scriptSource string
}

func NewScript(lang, script string) Query {
	return &Script{
		lang:         lang,
		scriptSource: script,
	}
}

func (s Script) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	response["lang"] = s.lang
	response["source"] = s.scriptSource
	return response
}

type TermsQuery struct {
	field string
	value interface{}
}

func NewTermsQuery(field string, value interface{}) Query {
	return &TermsQuery{
		field: field,
		value: value,
	}
}

func (s TermsQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	response["terms"] = map[string]interface{}{
		s.field: s.value,
	}
	return response
}

type MatchPhraseQuery struct {
	field string
	value string
	boost float64
}

func NewMatchPhraseQuery(field, value string, boost float64) Query {
	return &MatchPhraseQuery{
		field: field,
		value: value,
		boost: boost,
	}
}

func (s MatchPhraseQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	if s.boost == 0 {
		s.boost = 1
	}
	response["match_phrase"] = map[string]interface{}{
		s.field: map[string]interface{}{
			"query": s.value,
			"boost": s.boost,
		},
	}
	return response
}

type MatchQuery struct {
	field string
	value string
	boost float64
}

func NewMatchQuery(field, value string, boost float64) Query {
	return &MatchQuery{
		field: field,
		value: value,
		boost: boost,
	}
}

func (s MatchQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	if s.boost == 0 {
		s.boost = 1
	}
	response["match"] = map[string]interface{}{
		s.field: map[string]interface{}{
			"query": s.value,
			"boost": s.boost,
		},
	}
	return response
}

type BoolQuery struct {
	should Queries
	must   Queries
}

func NewBoolQuery(should, must Queries) Query {
	return &BoolQuery{
		must:   must,
		should: should,
	}
}

func (s BoolQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	var value = make(map[string]interface{})
	if len(s.should) > 0 {
		value["should"] = s.should.Source()
	}
	if len(s.must) > 0 {
		value["must"] = s.must.Source()
	}
	response["bool"] = value
	return response
}

type FunctionQuery struct {
	bool Query
}

func NewFunctionQuery(bool Query) Query {
	return &FunctionQuery{bool: bool}
}

func (s FunctionQuery) Source() map[string]interface{} {
	return s.bool.Source()
}

type Filter struct {
	filterType string
	field      string
	value      interface{}
}

func NewFilter(filterType, field string, value interface{}) Query {
	return &Filter{
		filterType: filterType,
		field:      field,
		value:      value,
	}
}

func (s Filter) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	if len(s.filterType) > 0 && len(s.field) > 0 {
		response[s.filterType] = map[string]interface{}{s.field: s.value}
	}
	return response
}

type FunctionFilter struct {
	filter Query
	weight float64
}

func NewFunctionFilter(weight float64, filter Query) Query {
	return &FunctionFilter{
		filter: filter,
		weight: weight,
	}
}

func (s FunctionFilter) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	response["weight"] = s.weight
	response["filter"] = s.filter.Source()
	return response
}

type FunctionScoreQuery struct {
	query     Query
	functions Queries
	boostMode string
}

func NewFunctionScoreQuery(boostMode string, query Query, queries Queries) Query {
	return &FunctionScoreQuery{
		boostMode: boostMode,
		query:     query,
		functions: queries,
	}
}

func (s FunctionScoreQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	if len(s.boostMode) == 0 {
		s.boostMode = "multiply"
	}

	response["function_score"] = map[string]interface{}{
		"query":      s.query.Source(),
		"boost_mode": s.boostMode,
		"functions":  s.functions.Source(),
	}
	return response
}

type MultiMatchQuery struct {
	query              string
	matchType          string
	analyzer           string
	fields             []string
	minimumShouldMatch string
}

func NewMultiMatchQuery(query, matchType, analyzer, minimumShouldMatch string, fields []string) Query {
	return &MultiMatchQuery{
		query:              query,
		matchType:          matchType,
		analyzer:           analyzer,
		minimumShouldMatch: minimumShouldMatch,
		fields:             fields,
	}
}

func (m MultiMatchQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	var query = make(map[string]interface{})
	query["query"] = m.query
	query["fields"] = m.fields
	if m.matchType == "" {
		m.matchType = "best_fields"
	}
	query["type"] = m.matchType
	if m.analyzer != "" {
		query["analyzer"] = m.analyzer
	}
	if m.minimumShouldMatch == "" {
		m.minimumShouldMatch = "100%"
	}
	query["minimum_should_match"] = m.minimumShouldMatch

	response["multi_match"] = query
	return response
}

type FuzzyQuery struct {
	Type           string
	Value          string
	Fuzziness      int
	PrefixLength   int
	Transpositions bool
	MaxExpansions  int
}

func NewFuzzyQuery(field, value string, fuzziness, prefixLength, maxExpansions int, transpositions bool) Query {
	return &FuzzyQuery{
		Type:           field,
		Value:          value,
		Fuzziness:      fuzziness,
		PrefixLength:   prefixLength,
		MaxExpansions:  maxExpansions,
		Transpositions: transpositions,
	}
}

func (f FuzzyQuery) Source() map[string]interface{} {
	var response = make(map[string]interface{})
	var queryType = make(map[string]interface{})
	var query = make(map[string]interface{})
	query["value"] = f.Value
	query["fuzziness"] = f.Fuzziness
	query["prefix_length"] = f.PrefixLength
	query["max_expansions"] = f.MaxExpansions
	query["transpositions"] = f.Transpositions
	queryType[f.Type] = query
	response["fuzzy"] = queryType
	return response
}
