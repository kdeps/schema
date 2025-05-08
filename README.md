# Kdeps Schema

This is the schema definitions used by [kdeps](https://kdeps.com).
See the [schema documentation](https://kdeps.github.io/schema).

## What is Kdeps?

Kdeps is an AI Agent framework for building self-hosted RAG AI Agents powered by open-source LLMs.

## Release Notes

### Latest Release: v0.2.24
  - Session sqlite storage (#14)
    * added deleteItem to memory storage
    
    * add Session ephemeral storage
    updated relnote/readme

### Previous Highlights
- **v0.2.23**:   - Expr block (#13)
    * removed unusued MemoryRecord class
    
    * add 'expr {...}' block to eval pkl expressions, i.e. memory.setItem('foo', 'bar')
    Merge branch 'main' of https://github.com/kdeps/schema
    
    updated relnote/readme

- **v0.2.22**:   - simplify memory operations (getItem/setItem), and add clear() (#12)

- **v0.2.21**:   - add the persistent sqlite memory record item read and update function (#11)
    
    Merge branch 'main' of https://github.com/kdeps/schema
    
    update readme/relnote

- **v0.2.20**:   - remove specific LLM role types (#10)
    
    updated schema readme/relnotes

- **v0.2.19**:   - Multi-prompt support (#9)
    * add role types on LLM action
    
    * add multiprompt support

- **v0.2.18**:   - add role types on LLM action (#8)

- **v0.2.17**:   - Add TrustedProxies to WebServer (#7)
    * fix indentation on deps/pkl files
    
    * add trustedproxy on webserver
    
    * keep webserver name in parity with apiserver

- **v0.2.16**:   - add webserver settings to project workflow settings (#6)

- **v0.2.15**:   - Create frontend serving settings (#5)
    * upgrade default ollama version to 0.6.5
    
    * added CORS configuration to APIServer
    
    * added webserver settings for serving static (html, htmx, etc.) or app (nodejs, streamlit, rails, etc.)

- **v0.2.14**:   - allow restricting http methods and params per resource (#4)
    
    Merge branch 'main' of https://github.com/kdeps/schema
    
    updated release notes

- **v0.2.13**:   - allow setting permitted HTTP values (headers, params, methods, routes) on resource (#3)
    
    updated release notes

- **v0.2.12**:   - add ability to set tz identifier timezone (#2)

- **v0.2.11**:   - Merge pull request #1 from kdeps/bump_versions_04_2025
    upgrade lowest pkl version to 0.28.1
    Merge branch 'bump_versions_04_2025' of https://github.com/kdeps/schema into bump_versions_04_2025
    
    upgrade lowest pkl version to 0.28.1
    
    upgrade lowest pkl version to 0.28.1

- **v0.2.10**:   - changed timestamp to duration from durationunit
    
    updated relnotes / readme

- **v0.2.9**:   - change timeoutDuration to Duration and use PKL semantics for duration seconds

- **v0.2.8**:   - Use DurationUnit for Timestamps; Upgrade pkl-go to 0.9.0

- **v0.2.7**:   - api response meta blocks can be optional
    
    updated release info

- **v0.2.6**:   - added api response meta info for other additional information sent over the JSON response
    
    updated release notes

- **v0.2.5**:   - allow sending headers to response

- **v0.2.4**:   - allow access to client IP and request ID
    
    added new schema documentation info
    
    updated release notes

- **v0.2.3**:   - added trustedProxies settings to API server

- **v0.2.2**:   - Change resource ID to actionID, and Workflow action to targetActionID

- **v0.2.1**:   - bump pkl to 0.27.2
    
    changed timeoutSeconds -> timeoutDuration

- **v0.2.0**:   - Use uniform naming convention for {Http,Api,Id,Json} -> {HTTP,API,ID,JSON}
    
    updated .gitattributes

- **v0.1.46**:   - reprioritize request skip validations

- **v0.1.45**:   - removed all deprecated imports

- **v0.1.44**:   - skip & validation is now a listing of any types

- **v0.1.43**:   - updated README.md
    
    bump pkl to 0.27.1

- **v0.1.42**:   - upgrade minimum pkl version to 0.26.3. added minor fixes on the relnote generator.
    
    Added skip validation helper functions

- **v0.1.41**:   - return stderr if not empty on stdout function
    
    added README.md and relnote generator

- **v0.1.40**:   - change request function from param("..") -> params("..")

- **v0.1.39**:   - import upstream PKL modules and KDEPS PKL helpers in resource & api response
    
    document renderers respond with null rather than empty string

- **v0.1.38**:   - added document pkl module for parsing and creating json, yaml and xml docs

- **v0.1.37**:   - added Data resource helper for getting agent data file path

- **v0.1.36**:   - Added file attribute where each associated resource value was saved

- **v0.1.35**:   - decode base64 strings by default on all Resource types

- **v0.1.34**:   - added params mapping to http client resources (go gen code)
    
    added params mapping to http client resources

- **v0.1.33**:   - make the build args as mapping; added build env in docker settings

- **v0.1.32**:   - Add parameters to docker settings; removed unused PKL configurations

- **v0.1.31**:   - added heroImage and agentIcon
    
    added ollamaImageTag property in workflow settings

- **v0.1.30**:   - Make the API response errors block listing (array)

- **v0.1.29**:   - added condaPackages section
    
    added support for installing conda packages

- **v0.1.28**:   - register python resource to resource pkl
    
    Revert "upgrade pkl to 0.27.0; register python resource to resource pkl"
    This reverts commit 86a334d697479e307513f759d2a7b0b06f9be35c.
    
    Revert "Update Gradle to 8.10.2"
    This reverts commit 45a93d3ebfa5aa2e795144984f0cde22bc1dc127.

- **v0.1.27**:   - Update Gradle to 8.10.2

- **v0.1.26**:   - upgrade pkl to 0.27.0; register python resource to resource pkl

- **v0.1.25**:   - added python script resource

- **v0.1.24**:   - reflect name changes to dockersettings go source

- **v0.1.23**:   - changed ppa to repositories
    
    added install Anaconda option

- **v0.1.22**:   - Make workflow an open module

- **v0.1.21**:   - deprecate postflightCheck

- **v0.1.20**:   - upgrade cicd pkl to 0.26.1

- **v0.1.19**:   - added API documentation. Change request API files functions for clarity

- **v0.1.18**:   - added Ubuntu PPA support; add API for querying filetypes by index

- **v0.1.17**:   - simplify usage of vision models and attachments

- **v0.1.16**:   - added visionFiles to LLM

- **v0.1.15**:   - fix typo for resolving filepath

- **v0.1.14**:   - renamed responseFile to file
    
    renamed responseFile to file
    
    added API server responseFile
    
    added request APIs for getting the list of files

- **v0.1.13**:   - always return the first element, in case of a single file upload
    
    support multiple file uploads

- **v0.1.12**:   - added vision and image gen attributes to LLM

- **v0.1.11**:   - make api serverr request extendable

- **v0.1.10**:   - Added request file operations. Decode base64 request by default. include allowed HTTP methods.

- **v0.1.9**:   - fix isEmpty on class

- **v0.1.8**:   - simplify resource resolver methods
    
    remove api server response type, focus on json response for now

- **v0.1.7**:   - removed textproto from api response types (gen code)
    
    removed textproto from api response types

- **v0.1.6**:   - response keys gen code

- **v0.1.5**:   - added json response keys to LLM

- **v0.1.4**:   - fix http response block

- **v0.1.3**:   - change schema to jsonResponse (bool)

- **v0.1.2**:   - set empty defaults when calling a function
    
    added security enforcement settings
    
    added schema to LLM

- **v0.1.1**:   - added extended functions for exec, http and llm modules

- **v0.1.0**:   - added deferred response api
    
    added llm timeout and kdeps dir settings
    
    add retry mechanism on failed create docs step

- **v0.0.50**:   - increase build timeout to 1 minute
    
    use a unified api for accessing resource values

- **v0.0.49**:   - exec, chat and client are all blocks

- **v0.0.48**:   - pre/post flights now a block that includes custom api error

- **v0.0.47**:   - change resource action to be not a listing

- **v0.0.46**:   - Revert "disable docs gha action until sorted this out"
    This reverts commit 6bed01569525deb1c8fc803b788c9391b1540b01.

- **v0.0.45**:   - disable docs gha action until sorted this out
    
    added apiResponse to resource; rename run -> action

- **v0.0.44**:   - simplify api server request template by using Dynamic maps

- **v0.0.43**:   - workflow and api server has a single required action

- **v0.0.42**:   - non-optional fields and use listing boolean for prefligh & skip steps

- **v0.0.41**:   - dockerSettings now becomes agentSettings since it has LLM models & hostIP/portNum was moved to apiSettings
    
    GEN: move the hostIP and serverPort settings to API Server settings

- **v0.0.40**:   - move the hostIP and serverPort settings to API Server settings

- **v0.0.38**:   - Added API Server Request/Response templates
    
    added Tags, a globally referenceable token

- **v0.0.37**:   - Changed docker hostIP portbinding settings

- **v0.0.36**:   - Changed docker settings for hostName and portNum; Fixed port default value

- **v0.0.35**:   - added hostname and portnum settings to docker

- **v0.0.34**:   - remove the placeholders for resource and workflow
    
    remove the placeholders for resource and workflow

- **v0.0.33**:   - Added semantics for external workflows; Additional semantics for resource dependencies

- **v0.0.32**:   - updated gen sdk

- **v0.0.31**:   - added optional params for default template fields on resource

- **v0.0.30**:   - Removed llm-apis (for now) and make local LLM the default. (reinstate commercial and cloud llm-apis in future versions)
    
    Removed workflows array, all pkl files within resources/ folder will be a workflow. Make a resource not an array but a single entry

- **v0.0.29**:   - fix workflow template validation

- **v0.0.28**:   - updated gen sdks
    
    removed RAG resource, and expect condition, renamed check to preflight, renamed api to httpclient
    
    removed modelfile, parameters, schema chat
    
    removed interactive input for ENV
    
    renamed ResourceAPI to ResourceHTTPClient

- **v0.0.27**:   - deprecate templates on build
    
    added docker container init settings

- **v0.0.26**:   - updated gen code with dockerGPU and runMode setings
    
    added dockerGPU and runMode to system config

- **v0.0.25**:   - 

- **v0.0.24**:   - updated kdeps template with new fields

- **v0.0.23**:   - reinstate read resource

- **v0.0.22**:   - init llmapikeys on settings

- **v0.0.21**:   - removed read resource for now

- **v0.0.20**:   - added default value for resource env read

- **v0.0.19**:   - read llm api keys from env vars

- **v0.0.18**:   - Updated kdeps.pkl to add global configs
    
    Added docker image and llmsettings directive

- **v0.0.17**:   - added template generation step

- **v0.0.16**:   - removed login ci step

- **v0.0.15**:   - upload pkl to firebase storage

- **v0.0.14**:   - regexp match for go version semantics - PklProject

- **v0.0.13**:   - regexp match for go version semantics - GHA

- **v0.0.12**:   - regexp match for go version semantics - GHA

- **v0.0.11**:   - regexp match for go version semantics

- **v0.0.2**:   - try adding prefixed 'v' for go mod tidy

- **v0.0.1**:   - 

- **0.0.10**:   - init go.mod; move schema/pkg to schema/gen
    
    added sleep timer to resolve new pkl version release

- **0.0.9**:   - add global settings for kdeps
    
    initialize the default routes
    
    renamed Settings.pkl -> Project.pkl

- **0.0.8**:   - Added API server mode
    
    fix indentation

- **0.0.7**:   - added project support files
    
    move schema module path to deps/pkl

- **0.0.6**:   - specify the source url pattern in Pkl docs

- **0.0.5**:   - use pkl-docs env for deploying pages
    
    Added Makefile and gen-go code
    
    Delete CNAME
    
    Update CNAME
    
    Create CNAME

- **0.0.4**:   - Fix resolution of pkl files

- **0.0.3**:   - Fix zip package resolution in PklProject; Use pkg as gen folder

- **0.0.2**:   - fix version resolution

