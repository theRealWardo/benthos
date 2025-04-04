package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"gopkg.in/yaml.v3"

	"github.com/benthosdev/benthos/v4/internal/bloblang/query"
	"github.com/benthosdev/benthos/v4/internal/bundle"
	"github.com/benthosdev/benthos/v4/internal/docs"
	"github.com/benthosdev/benthos/v4/internal/manager"
)

// ConfigField describes a field within a component configuration, to be added
// to a ConfigSpec.
type ConfigField struct {
	field docs.FieldSpec
}

// NewAnyField describes a new config field that can assume any value type
// without triggering a config parse or linting error.
func NewAnyField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldAnything(name, ""),
	}
}

// NewAnyListField describes a new config field consisting of a list of values
// that can assume any value type without triggering a config parse or linting
// error.
func NewAnyListField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldAnything(name, "").Array(),
	}
}

// NewAnyMapField describes a new config field consisting of a map of values
// that can assume any value type without triggering a config parse or linting
// error.
func NewAnyMapField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldAnything(name, "").Map(),
	}
}

// NewStringField describes a new string type config field.
func NewStringField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldString(name, ""),
	}
}

// NewDurationField describes a new duration string type config field, allowing
// users to define a time interval with strings of the form 60s, 3m, etc.
func NewDurationField(name string) *ConfigField {
	// TODO: Add linting rule for duration
	return &ConfigField{
		field: docs.FieldString(name, ""),
	}
}

// NewStringEnumField describes a new string type config field that can have one
// of a discrete list of values.
func NewStringEnumField(name string, options ...string) *ConfigField {
	return &ConfigField{
		field: docs.FieldString(name, "").HasOptions(options...),
	}
}

// NewStringAnnotatedEnumField describes a new string type config field that can
// have one of a discrete list of values, where each value must be accompanied
// by a description that annotates its behaviour in the documentation.
func NewStringAnnotatedEnumField(name string, options map[string]string) *ConfigField {
	optionKeys := make([]string, 0, len(options))
	for key := range options {
		optionKeys = append(optionKeys, key)
	}
	sort.Strings(optionKeys)

	flatOptions := make([]string, 0, len(options)*2)
	for _, o := range optionKeys {
		flatOptions = append(flatOptions, o, options[o])
	}

	return &ConfigField{
		field: docs.FieldString(name, "").HasAnnotatedOptions(flatOptions...),
	}
}

// NewStringListField describes a new config field consisting of a list of
// strings.
func NewStringListField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldString(name, "").Array(),
	}
}

// NewStringMapField describes a new config field consisting of an object of
// arbitrary keys with string values.
func NewStringMapField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldString(name, "").Map(),
	}
}

// NewIntField describes a new int type config field.
func NewIntField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldInt(name, ""),
	}
}

// NewIntListField describes a new config field consisting of a list of
// integers.
func NewIntListField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldInt(name, "").Array(),
	}
}

// NewIntMapField describes a new config field consisting of an object of
// arbitrary keys with integer values.
func NewIntMapField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldInt(name, "").Map(),
	}
}

// NewFloatField describes a new float type config field.
func NewFloatField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldFloat(name, ""),
	}
}

// NewBoolField describes a new bool type config field.
func NewBoolField(name string) *ConfigField {
	return &ConfigField{
		field: docs.FieldBool(name, ""),
	}
}

// NewObjectField describes a new object type config field, consisting of one
// or more child fields.
func NewObjectField(name string, fields ...*ConfigField) *ConfigField {
	children := make([]docs.FieldSpec, len(fields))
	for i, f := range fields {
		children[i] = f.field
	}
	return &ConfigField{
		field: docs.FieldObject(name, "").WithChildren(children...),
	}
}

// NewObjectListField describes a new list type config field consisting of
// objects with one or more child fields.
func NewObjectListField(name string, fields ...*ConfigField) *ConfigField {
	objField := NewObjectField(name, fields...)
	return &ConfigField{
		field: objField.field.Array(),
	}
}

// NewInternalField returns a ConfigField derived from an internal package field
// spec. This function is for internal use only.
func NewInternalField(ifield docs.FieldSpec) *ConfigField {
	return &ConfigField{
		field: ifield,
	}
}

// Description adds a description to the field which will be shown when printing
// documentation for the component config spec.
func (c *ConfigField) Description(d string) *ConfigField {
	c.field.Description = d
	return c
}

// Advanced marks a config field as being advanced, and therefore it will not
// appear in simplified documentation examples.
func (c *ConfigField) Advanced() *ConfigField {
	c.field = c.field.Advanced()
	return c
}

// Deprecated marks a config field as being deprecated, and therefore it will not
// appear in documentation examples.
func (c *ConfigField) Deprecated() *ConfigField {
	c.field = c.field.Deprecated()
	return c
}

// Default specifies a default value that this field will assume if it is
// omitted from a provided config. Fields that do not have a default value are
// considered mandatory, and so parsing a config will fail in their absence.
func (c *ConfigField) Default(v any) *ConfigField {
	c.field = c.field.HasDefault(v)
	return c
}

// Optional specifies that a field is optional even when a default value has not
// been specified. When a field is marked as optional you can test its presence
// within a parsed config with the method Contains.
func (c *ConfigField) Optional() *ConfigField {
	c.field = c.field.Optional()
	return c
}

// Secret marks this field as being a secret, which means it represents
// information that is generally considered sensitive such as passwords or
// access tokens.
func (c *ConfigField) Secret() *ConfigField {
	c.field = c.field.Secret()
	return c
}

// Example adds an example value to the field which will be shown when printing
// documentation for the component config spec.
func (c *ConfigField) Example(e any) *ConfigField {
	c.field.Examples = append(c.field.Examples, e)
	return c
}

// Version specifies the specific version at which this field was added to the
// component.
func (c *ConfigField) Version(v string) *ConfigField {
	c.field = c.field.AtVersion(v)
	return c
}

// LintRule adds a custom linting rule to the field in the form of a bloblang
// mapping. The mapping is provided the value of the field within a config as
// the context `this`, and if the mapping assigns to `root` an array of one or
// more strings these strings will be exposed to a config author as linting
// errors.
//
// For example, if we wanted to add a linting rule for a string field that
// ensures the value contains only lowercase values we might add the following
// linting rule:
//
// `root = if this.lowercase() != this { [ "field must be lowercase" ] }`.
func (c *ConfigField) LintRule(blobl string) *ConfigField {
	c.field = c.field.LinterBlobl(blobl)
	return c
}

//------------------------------------------------------------------------------

// ConfigSpec describes the configuration specification for a plugin
// component. This will be used for validating and linting configuration files
// and providing a parsed configuration struct to the plugin constructor.
type ConfigSpec struct {
	component docs.ComponentSpec
}

func (c *ConfigSpec) configFromNode(mgr bundle.NewManagement, node *yaml.Node) (*ParsedConfig, error) {
	fields, err := c.component.Config.YAMLToValue(node, docs.ToValueConfig{})
	if err != nil {
		return nil, err
	}

	return &ParsedConfig{mgr: mgr, generic: fields}, nil
}

// ParseYAML attempts to parse a YAML document as the defined configuration spec
// and returns a parsed config view. The provided environment determines which
// child components and Bloblang functions can be created by the fields of the
// spec, you can leave the environment nil to use the global environment.
//
// This method is intended for testing purposes and is not required for normal
// use of plugin components, as parsing is managed by other components.
func (c *ConfigSpec) ParseYAML(yamlStr string, env *Environment) (*ParsedConfig, error) {
	if env == nil {
		env = globalEnvironment
	}

	var nconf yaml.Node
	if err := yaml.Unmarshal([]byte(yamlStr), &nconf); err != nil {
		return nil, err
	}
	if nconf.Kind == yaml.DocumentNode && len(nconf.Content) > 0 {
		nconf = *nconf.Content[0]
	}

	mgr, err := manager.New(
		manager.NewResourceConfig(),
		manager.OptSetEnvironment(env.internal),
		manager.OptSetBloblangEnvironment(env.getBloblangParserEnv()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate resources: %w", err)
	}

	return c.configFromNode(mgr, &nconf)
}

// NewConfigSpec creates a new empty component configuration spec. If the
// plugin does not require configuration fields the result of this call is
// enough.
func NewConfigSpec() *ConfigSpec {
	return &ConfigSpec{
		component: docs.ComponentSpec{
			Status: docs.StatusExperimental,
			Plugin: true,
			Config: docs.FieldComponent(),
		},
	}
}

// Stable sets a documentation label on the component indicating that its
// configuration spec is stable. Plugins are considered experimental by default.
func (c *ConfigSpec) Stable() *ConfigSpec {
	c.component.Status = docs.StatusStable
	return c
}

// Beta sets a documentation label on the component indicating that its
// configuration spec is ready for beta testing, meaning backwards incompatible
// changes will not be made unless a fundamental problem is found. Plugins are
// considered experimental by default.
func (c *ConfigSpec) Beta() *ConfigSpec {
	c.component.Status = docs.StatusBeta
	return c
}

// Deprecated sets a documentation label on the component indicating that it is
// now deprecated. Plugins are considered experimental by default.
func (c *ConfigSpec) Deprecated() *ConfigSpec {
	c.component.Status = docs.StatusDeprecated
	return c
}

// Categories adds one or more string tags to the component, these are used for
// arbitrarily grouping components in documentation.
func (c *ConfigSpec) Categories(categories ...string) *ConfigSpec {
	c.component.Categories = categories
	return c
}

// Version specifies that this component was introduced in a given version.
func (c *ConfigSpec) Version(v string) *ConfigSpec {
	c.component.Version = v
	return c
}

// Summary adds a short summary to the plugin configuration spec that describes
// the general purpose of the component.
func (c *ConfigSpec) Summary(summary string) *ConfigSpec {
	c.component.Summary = summary
	return c
}

// Description adds a description to the plugin configuration spec that
// describes in more detail the behaviour of the component and how it should be
// used.
func (c *ConfigSpec) Description(description string) *ConfigSpec {
	c.component.Description = description
	return c
}

// Field sets the specification of a field within the config spec, used for
// linting and generating documentation for the component.
//
// If the provided field has an empty name then it registered as the value at
// the root of the config spec.
//
// When creating a spec with a struct constructor the fields from that struct
// will already be inferred. However, setting a field explicitly is sometimes
// useful for enriching the field documentation with more information.
func (c *ConfigSpec) Field(f *ConfigField) *ConfigSpec {
	if f.field.Name == "" {
		// Set field to root of config spec
		c.component.Config = f.field
		return c
	}

	c.component.Config.Type = docs.FieldTypeObject
	for i, s := range c.component.Config.Children {
		if s.Name == f.field.Name {
			c.component.Config.Children[i] = f.field
			return c
		}
	}

	c.component.Config.Children = append(c.component.Config.Children, f.field)
	return c
}

// Example adds an example to the plugin configuration spec that demonstrates
// how the component can be used. An example has a title, summary, and a YAML
// configuration showing a real use case.
func (c *ConfigSpec) Example(title, summary, config string) *ConfigSpec {
	c.component.Examples = append(c.component.Examples, docs.AnnotatedExample{
		Title:   title,
		Summary: summary,
		Config:  config,
	})
	return c
}

// EncodeJSON attempts to parse a JSON object as a byte slice and uses it to
// populate the configuration spec. The schema of this method is undocumented
// and is not intended for general use.
//
// Experimental: This method is not intended for general use and could have its
// signature and/or behaviour changed outside of major version bumps.
func (c *ConfigSpec) EncodeJSON(v []byte) error {
	return json.Unmarshal(v, &c.component)
}

// LintRule adds a custom linting rule to the ConfigSpec in the form of a
// bloblang mapping. The mapping is provided the value of the fields within
// the ConfigSpec as the context `this`, and if the mapping assigns to `root` an
// array of one or more strings these strings will be exposed to a config author
// as linting errors.
//
// For example, if we wanted to add a linting rule for several ConfigSpec fields
// that ensures some fields are mutually exclusive and some require others we
// might use the following:
//
// root = match {
// this.exists("meow") && this.exists("woof") => [ "both `+"`meow`"+` and `+"`woof`"+` can't be set simultaneously" ],
// this.exists("reticulation") && (!this.exists("splines") || this.splines == "") => [ "`+"`splines`"+` is required when setting `+"`reticulation`"+`" ],
// }.
func (c *ConfigSpec) LintRule(blobl string) *ConfigSpec {
	c.component.Config = c.component.Config.LinterBlobl(blobl)
	return c
}

//------------------------------------------------------------------------------

// ConfigView is a struct returned by a Benthos service environment when walking
// the list of registered components and provides access to information about
// the component.
type ConfigView struct {
	component docs.ComponentSpec
}

// Summary returns a documentation summary of the component, often formatted as
// markdown.
func (c *ConfigView) Summary() string {
	return c.component.Summary
}

// Description returns a documentation description of the component, often
// formatted as markdown.
func (c *ConfigView) Description() string {
	return c.component.Description
}

// IsDeprecated returns true if the component is marked as deprecated.
func (c *ConfigView) IsDeprecated() bool {
	return c.component.Status == docs.StatusDeprecated
}

// FormatJSON returns a byte slice of the component configuration formatted as a
// JSON object. The schema of this method is undocumented and is not intended
// for general use.
//
// Experimental: This method is not intended for general use and could have its
// signature and/or behaviour changed outside of major version bumps.
func (c *ConfigView) FormatJSON() ([]byte, error) {
	return json.Marshal(c.component)
}

// RenderDocs creates a markdown file that documents the configuration of the
// component config view. This markdown may include Docusaurus react elements as
// it matches the documentation generated for the official Benthos website.
//
// Experimental: This method is not intended for general use and could have its
// signature and/or behaviour changed outside of major version bumps.
func (c *ConfigView) RenderDocs() ([]byte, error) {
	_, rootOnly := map[string]struct{}{
		"cache":      {},
		"rate_limit": {},
		"processor":  {},
	}[string(c.component.Type)]

	conf := map[string]any{
		"type": c.component.Name,
	}
	for k, v := range docs.ReservedFieldsByType(c.component.Type) {
		if k == "plugin" {
			continue
		}
		if v.Default != nil {
			conf[k] = *v.Default
		}
	}

	return c.component.AsMarkdown(!rootOnly, conf)
}

//------------------------------------------------------------------------------

// ParsedConfig represents a plugin configuration that has been validated and
// parsed from a ConfigSpec, and allows plugin constructors to access
// configuration fields.
type ParsedConfig struct {
	hiddenPath []string
	mgr        bundle.NewManagement
	generic    any
}

// Namespace returns a version of the parsed config at a given field namespace.
// This is useful for extracting multiple fields under the same grouping.
func (p *ParsedConfig) Namespace(path ...string) *ParsedConfig {
	tmpConfig := *p
	tmpConfig.hiddenPath = append([]string{}, p.hiddenPath...)
	tmpConfig.hiddenPath = append(tmpConfig.hiddenPath, path...)
	tmpConfig.mgr = p.mgr.IntoPath(path...)
	return &tmpConfig
}

// Field accesses a field from the parsed config by its name and returns the
// value if the field is found and a boolean indicating whether it was found.
// Nested fields can be accessed by specifying the series of field names.
func (p *ParsedConfig) field(path ...string) (any, bool) {
	gObj := gabs.Wrap(p.generic).S(p.hiddenPath...)
	if exists := gObj.Exists(path...); !exists {
		return nil, false
	}
	return gObj.S(path...).Data(), true
}

func (p *ParsedConfig) fullDotPath(path ...string) string {
	var fullPath []string
	fullPath = append(fullPath, p.hiddenPath...)
	fullPath = append(fullPath, path...)
	return strings.Join(fullPath, ".")
}

// Contains checks whether the parsed config contains a given field identified
// by its name.
func (p *ParsedConfig) Contains(path ...string) bool {
	gObj := gabs.Wrap(p.generic).S(p.hiddenPath...)
	return gObj.Exists(path...)
}

// FieldAny accesses a field from the parsed config by its name that can assume
// any value type. If the field is not found an error is returned.
func (p *ParsedConfig) FieldAny(path ...string) (any, error) {
	v, exists := p.field(path...)
	if !exists {
		return "", fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	return v, nil
}

// FieldAnyList accesses a field that is a list of any value types from the
// parsed config by its name and returns the value as an array of *ParsedConfig
// types, where each one represents an object or value in the list. Returns an
// error if the field is not found, or is not a list of values.
func (p *ParsedConfig) FieldAnyList(path ...string) ([]*ParsedConfig, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iList, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected field '%v' to be a list, got %T", p.fullDotPath(path...), v)
	}
	sList := make([]*ParsedConfig, len(iList))
	for i, ev := range iList {
		sList[i] = &ParsedConfig{
			mgr:     p.mgr,
			generic: ev,
		}
	}
	return sList, nil
}

// FieldAnyMap accesses a field that is an object of arbitrary keys and any
// values from the parsed config by its name and returns a map of *ParsedConfig
// types, where each one represents an object or value in the map. Returns an
// error if the field is not found, or is not an object.
func (p *ParsedConfig) FieldAnyMap(path ...string) (map[string]*ParsedConfig, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iMap, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected field '%v' to be a string map, got %T", p.fullDotPath(path...), v)
	}
	sMap := make(map[string]*ParsedConfig, len(iMap))
	for k, v := range iMap {
		sMap[k] = &ParsedConfig{
			mgr:     p.mgr,
			generic: v,
		}
	}
	return sMap, nil
}

// FieldString accesses a string field from the parsed config by its name. If
// the field is not found or is not a string an error is returned.
func (p *ParsedConfig) FieldString(path ...string) (string, error) {
	v, exists := p.field(path...)
	if !exists {
		return "", fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	str, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("expected field '%v' to be a string, got %T", p.fullDotPath(path...), v)
	}
	return str, nil
}

// FieldDuration accesses a duration string field from the parsed config by its
// name. If the field is not found or is not a valid duration string an error is
// returned.
func (p *ParsedConfig) FieldDuration(path ...string) (time.Duration, error) {
	v, exists := p.field(path...)
	if !exists {
		return 0, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	str, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("expected field '%v' to be a string, got %T", p.fullDotPath(path...), v)
	}
	d, err := time.ParseDuration(str)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%v' as a duration string: %w", p.fullDotPath(path...), err)
	}
	return d, nil
}

// FieldStringList accesses a field that is a list of strings from the parsed
// config by its name and returns the value. Returns an error if the field is
// not found, or is not a list of strings.
func (p *ParsedConfig) FieldStringList(path ...string) ([]string, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iList, ok := v.([]any)
	if !ok {
		if sList, ok := v.([]string); ok {
			return sList, nil
		}
		return nil, fmt.Errorf("expected field '%v' to be a string list, got %T", p.fullDotPath(path...), v)
	}
	sList := make([]string, len(iList))
	for i, ev := range iList {
		if sList[i], ok = ev.(string); !ok {
			return nil, fmt.Errorf("expected field '%v' to be a string list, found an element of type %T", p.fullDotPath(path...), ev)
		}
	}
	return sList, nil
}

// FieldStringMap accesses a field that is an object of arbitrary keys and
// string values from the parsed config by its name and returns the value.
// Returns an error if the field is not found, or is not an object of strings.
func (p *ParsedConfig) FieldStringMap(path ...string) (map[string]string, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iMap, ok := v.(map[string]any)
	if !ok {
		if sMap, ok := v.(map[string]string); ok {
			return sMap, nil
		}
		return nil, fmt.Errorf("expected field '%v' to be a string map, got %T", p.fullDotPath(path...), v)
	}
	sMap := make(map[string]string, len(iMap))
	for k, ev := range iMap {
		if sMap[k], ok = ev.(string); !ok {
			return nil, fmt.Errorf("expected field '%v' to be a string map, found an element of type %T", p.fullDotPath(path...), ev)
		}
	}
	return sMap, nil
}

// FieldInt accesses an int field from the parsed config by its name and returns
// the value. Returns an error if the field is not found or is not an int.
func (p *ParsedConfig) FieldInt(path ...string) (int, error) {
	v, exists := p.field(path...)
	if !exists {
		return 0, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	i, err := query.IGetInt(v)
	if err != nil {
		return 0, fmt.Errorf("expected field '%v' to be an int, got %T", p.fullDotPath(path...), v)
	}
	return int(i), nil
}

// FieldIntList accesses a field that is a list of integers from the parsed
// config by its name and returns the value. Returns an error if the field is
// not found, or is not a list of integers.
func (p *ParsedConfig) FieldIntList(path ...string) ([]int, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iList, ok := v.([]any)
	if !ok {
		if sList, ok := v.([]int); ok {
			return sList, nil
		}
		return nil, fmt.Errorf("expected field '%v' to be an integer list, got %T", p.fullDotPath(path...), v)
	}
	sList := make([]int, len(iList))
	for i, ev := range iList {
		iv, err := query.IToInt(ev)
		if err != nil {
			return nil, fmt.Errorf("expected field '%v' to be an integer list, found an element of type %T", p.fullDotPath(path...), ev)
		}
		sList[i] = int(iv)
	}
	return sList, nil
}

// FieldIntMap accesses a field that is an object of arbitrary keys and
// integer values from the parsed config by its name and returns the value.
// Returns an error if the field is not found, or is not an object of integers.
func (p *ParsedConfig) FieldIntMap(path ...string) (map[string]int, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iMap, ok := v.(map[string]any)
	if !ok {
		if sMap, ok := v.(map[string]int); ok {
			return sMap, nil
		}
		return nil, fmt.Errorf("expected field '%v' to be an integer map, got %T", p.fullDotPath(path...), v)
	}
	sMap := make(map[string]int, len(iMap))
	for k, ev := range iMap {
		iv, err := query.IToInt(ev)
		if err != nil {
			return nil, fmt.Errorf("expected field '%v' to be an integer map, found an element of type %T", p.fullDotPath(path...), ev)
		}
		sMap[k] = int(iv)
	}
	return sMap, nil
}

// FieldFloat accesses a float field from the parsed config by its name and
// returns the value. Returns an error if the field is not found or is not a
// float.
func (p *ParsedConfig) FieldFloat(path ...string) (float64, error) {
	v, exists := p.field(path...)
	if !exists {
		return 0, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	f, err := query.IGetNumber(v)
	if err != nil {
		return 0, fmt.Errorf("expected field '%v' to be a float, got %T", p.fullDotPath(path...), v)
	}
	return f, nil
}

// FieldBool accesses a bool field from the parsed config by its name and
// returns the value. Returns an error if the field is not found or is not a
// bool.
func (p *ParsedConfig) FieldBool(path ...string) (bool, error) {
	v, e := p.field(path...)
	if !e {
		return false, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("expected field '%v' to be a bool, got %T", p.fullDotPath(path...), v)
	}
	return b, nil
}

// FieldObjectList accesses a field that is a list of objects from the parsed
// config by its name and returns the value as an array of *ParsedConfig types,
// where each one represents an object in the list. Returns an error if the
// field is not found, or is not a list of objects.
func (p *ParsedConfig) FieldObjectList(path ...string) ([]*ParsedConfig, error) {
	v, exists := p.field(path...)
	if !exists {
		return nil, fmt.Errorf("field '%v' was not found in the config", p.fullDotPath(path...))
	}
	iList, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected field '%v' to be a list, got %T", p.fullDotPath(path...), v)
	}
	sList := make([]*ParsedConfig, len(iList))
	for i, ev := range iList {
		sList[i] = &ParsedConfig{
			mgr:     p.mgr,
			generic: ev,
		}
	}
	return sList, nil
}
