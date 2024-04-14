Commands | Explanation
---------|------------
compile  | This command helps check SQL syntax and reports any typing errors.
completion | This command is to generate an auto-completion script for your environment. The following are the supported environments: Bash, Fish, PowerShell, and zsh.
generate | A command to generate the .go files based on the provided SQL statements. This will be the command that we will be using a lot for the application.
init | This command is the first command that is used to initialize your application to start using this tool.

**Fields in yaml file**

Tag Name | Description
---------|------------
Name | Any string to be used as the package name.
Path | Specifies the name of the directory that will host the generated .go code.
Queries | Specifies the directory name containing the SQL queries that sqlc will use to generate the .go code.
Schema | A directory containing SQL files that will be used to generate all the relevant .go files.
Engine | Specifies the database engine that will be used: sqlc supports either MySQL or Postgres.
emit_db_tags | Setting this to true will generate the struct with db tags
emit_prepared_queries | Setting this to true instructs sqlc to support prepared queries in the generated code.
emit_interface | Setting this to true will instruct sqlc to generate the querier interface.
emit_exact_table_names | Setting this to true will instruct sqlc to mirror the struct name to the table name.
emit_empty_slices | Setting this to true will instruct sqlc to return an empty slice for returning data on many sides of the table.
emit_exported_queries| Setting this to true will instruct sqlc to allow the SQL statement used in the auto-generated code to be accessed by an outside package.
emit_json_tags | Setting this to true will generate the struct with JSON tags.
json_tags_case_style | This setting can accept the following – camel, pascal, snake, and none. The case style is used for the JSON tags used in the struct. Normally, this is used with emit_json_tags.
output_db_file_name | Name used as the filename for the auto-generated database file.
output_models_file_name | Name used as the filename for the auto-generated model file.
output_querier_file_name | Name used as the filename for the auto-generated querier file.
output_files_suffix | Suffix to be used as part of the auto-generated query file.