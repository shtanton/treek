# treek

Like awk but for trees.

Currently supports only JSON.

Currently implemented in go but once the spec is final I'll reimplement in C or something.

# Examples

#### Extract a value
```
treek "path.to.value"
```

#### Print full names of all people
```
treek 'people.* {println($0.first_name + " " + $0.last_name)}'
```

#### Print average age of all people
```
treek 'people.*.age {total += $0; count += 1} {println(total / count)}'
```

#### Print the first names of all the Johnsons
```
treek 'people.($0.last_name=="Johnson").first_name'
```
