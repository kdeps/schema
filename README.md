# Kdeps Schema

This is the schema definitions used by [kdeps](https://kdeps.com).
See the [schema documentation](https://kdeps.github.io/schema).

## What is Kdeps?

Kdeps is an AI Agent framework for building self-hosted RAG AI Agents powered by open-source LLMs.

## Release Notes

### Latest Release: v0.3.2
*Released: 2025-07-04*


✨ **Enhancements**
  - **added pklresource and add use sensible defaults to optional fields (#20)** (`1496555`)

📦 **Updates**
  - **Update release notes for v0.3.1 [skip ci]** (`151b018`)

### Complete Release History

*Detailed changelog showing all changes from the beginning of the project*


## v0.3.1 (*2025-07-02*)

📦 **Updates**
  - **upgrade deprecated gha plugin** (`29f15fa`)

📝 **Other Changes**
  - **renamed Command->Script and CondaEnvironment->PythonEnvironment** (`b7d3688`)
  - **install pkl-gen-go in gha for running pkl tests** (`aaadde8`)
  - **install pkl-gen-go in gha for running pkl tests** (`5317aa2`)
  - **fix gha build** (`650674f`)
  - **enhance workflows to gen relnotes on gha trigger** (`f088834`)

## v0.3.0 (*2025-07-02*)

📝 **Other Changes**
  - **Uniform schema with Retry Logic and Docker Enhancements (#19)** (`64f620e`)
    * use uniform schema naming conventions (attr as capitalized, func as pascalCased)
  - **** (`* uniform schema`)

## v0.2.40 (*2025-06-11*)

📝 **Other Changes**
  - **embed deps/pkl schema files via go:embed** (`159f785`)

## v0.2.39 (*2025-06-11*)

✨ **Enhancements**
  - **add the pkl files in the output_dir for expose to tests** (`72b3641`)

## v0.2.38 (*2025-05-30*)

✨ **Enhancements**
  - **add resource#itemValues function to obtain item iteration values** (`9c48aa0`)

📦 **Updates**
  - **update readme/relnotes** (`7775ce8`)

## v0.2.37 (*2025-05-30*)

✨ **Enhancements**
  - **add pkl:json imports per each resource** (`8b252f4`)

## v0.2.36 (*2025-05-30*)

✨ **Enhancements**
  - **added the itemValues per resource to obtain the iteration results** (`59cadd4`)

## v0.2.35 (*2025-05-30*)

📝 **Other Changes**
  - **null propagate the array obtained from dynamic reader for item values** (`282f277`)

## v0.2.34 (*2025-05-30*)

✨ **Enhancements**
  - **add nullable defaults for item values listing** (`ef3bfc8`)

## v0.2.33 (*2025-05-30*)

✨ **Enhancements**
  - **add return value listing type for item** (`0b0f475`)

## v0.2.32 (*2025-05-30*)

✨ **Enhancements**
  - **return a new listing when item results is null** (`e8a979b`)

## v0.2.31 (*2025-05-28*)

📝 **Other Changes**
  - **values now require to pass the actionID (#18)** (`1cc199d`)

## v0.2.30 (*2025-05-24*)

📦 **Updates**
  - **update readme/relnotes** (`6fdaccb`)

📝 **Other Changes**
  - **removed id params on item operations** (`c035c8a`)

## v0.2.29 (*2025-05-24*)

📝 **Other Changes**
  - **change item function signature to not require an id params (#17)** (`034ac3e`)
    * added ability to iterate through items
  - **** (`* changed item.fetch -> item.current`)
  - **** (`* change item function signature to not require an id params`)

## v0.2.28 (*2025-05-23*)

✨ **Enhancements**
  - **added ability to iterate through items (#16)** (`c46cc9c`)
    * added ability to iterate through items

📦 **Updates**
  - **update readme/relnotes** (`c277d2a`)

📝 **Other Changes**
  - **** (`* changed item.fetch -> item.current`)

## v0.2.27 (*2025-05-12*)

📦 **Updates**
  - **updated relnote/readme** (`37fc99a`)

📝 **Other Changes**
  - **hotfix: register tool into resources** (`6d844d6`)

## v0.2.26 (*2025-05-12*)

✨ **Enhancements**
  - **hotfix: add path to the script or inline script to LLM tools** (`f461f98`)

## v0.2.25 (*2025-05-12*)

✨ **Enhancements**
  - **add ability for LLM to use tools (akin to MCP) (#15)** (`d49777f`)

📦 **Updates**
  - **updated relnote/readme** (`cb21d6a`)

## v0.2.24 (*2025-05-08*)

📦 **Updates**
  - **updated relnote/readme** (`5848573`)

📝 **Other Changes**
  - **Session sqlite storage (#14)** (`cf4b6b1`)
    * added deleteItem to memory storage
  - **** (`* add Session ephemeral storage`)

## v0.2.23 (*2025-05-08*)

📦 **Updates**
  - **updated relnote/readme** (`12ac9e5`)

📝 **Other Changes**
  - **Expr block (#13)** (`864ae2e`)
    * removed unusued MemoryRecord class
  - **** (`* add 'expr {...}' block to eval pkl expressions, i.e. memory.setItem('foo', 'bar')`)

## v0.2.22 (*2025-05-07*)

✨ **Enhancements**
  - **simplify memory operations (getItem/setItem), and add clear() (#12)** (`187e6b9`)

## v0.2.21 (*2025-05-06*)

✨ **Enhancements**
  - **add the persistent sqlite memory record item read and update function (#11)** (`0131739`)

📦 **Updates**
  - **update readme/relnote** (`6e92253`)

## v0.2.20 (*2025-04-28*)

📦 **Updates**
  - **updated schema readme/relnotes** (`27d3e8f`)

📝 **Other Changes**
  - **remove specific LLM role types (#10)** (`868bac5`)

## v0.2.19 (*2025-04-28*)

📝 **Other Changes**
  - **Multi-prompt support (#9)** (`cd4b60a`)
    * add role types on LLM action
  - **** (`* add multiprompt support`)

## v0.2.18 (*2025-04-26*)

✨ **Enhancements**
  - **add role types on LLM action (#8)** (`638e38d`)

## v0.2.17 (*2025-04-18*)

✨ **Enhancements**
  - **Add TrustedProxies to WebServer (#7)** (`f9e658f`)
    * fix indentation on deps/pkl files

📝 **Other Changes**
  - **** (`* add trustedproxy on webserver`)
  - **** (`* keep webserver name in parity with apiserver`)

## v0.2.16 (*2025-04-17*)

✨ **Enhancements**
  - **add webserver settings to project workflow settings (#6)** (`01fa2e1`)

## v0.2.15 (*2025-04-17*)

📝 **Other Changes**
  - **Create frontend serving settings (#5)** (`c7f1374`)
    * upgrade default ollama version to 0.6.5
  - **** (`* added CORS configuration to APIServer`)
  - **** (`* added webserver settings for serving static (html, htmx, etc.) or app (nodejs, streamlit, rails, etc.)`)

## v0.2.14 (*2025-04-16*)

📦 **Updates**
  - **updated release notes** (`b6bc1bd`)

📝 **Other Changes**
  - **allow restricting http methods and params per resource (#4)** (`329c08a`)

## v0.2.13 (*2025-04-16*)

📦 **Updates**
  - **updated release notes** (`fcc686b`)

📝 **Other Changes**
  - **allow setting permitted HTTP values (headers, params, methods, routes) on resource (#3)** (`011d1f4`)

## v0.2.12 (*2025-04-16*)

✨ **Enhancements**
  - **add ability to set tz identifier timezone (#2)** (`9a6fa07`)

## v0.2.11 (*2025-04-15*)

📦 **Updates**
  - **upgrade lowest pkl version to 0.28.1** (`bb44da6`)
  - **upgrade lowest pkl version to 0.28.1** (`b74f5a6`)

## v0.2.10 (*2025-02-16*)

📦 **Updates**
  - **updated relnotes / readme** (`6cc68de`)

📝 **Other Changes**
  - **changed timestamp to duration from durationunit** (`c4281eb`)

## v0.2.9 (*2025-02-16*)

📝 **Other Changes**
  - **change timeoutDuration to Duration and use PKL semantics for duration seconds** (`d4bb52a`)

## v0.2.8 (*2025-02-16*)

📦 **Updates**
  - **Use DurationUnit for Timestamps; Upgrade pkl-go to 0.9.0** (`d22127a`)

## v0.2.7 (*2025-02-10*)

📦 **Updates**
  - **updated release info** (`2d17f13`)

📝 **Other Changes**
  - **api response meta blocks can be optional** (`6b4869f`)

## v0.2.6 (*2025-02-10*)

✨ **Enhancements**
  - **added api response meta info for other additional information sent over the JSON response** (`cf98a9d`)

📦 **Updates**
  - **updated release notes** (`c085ef2`)

## v0.2.5 (*2025-02-10*)

📝 **Other Changes**
  - **allow sending headers to response** (`12f3575`)

## v0.2.4 (*2025-02-08*)

✨ **Enhancements**
  - **added new schema documentation info** (`c3fc856`)

📦 **Updates**
  - **updated release notes** (`0e6e41a`)

📝 **Other Changes**
  - **allow access to client IP and request ID** (`34a6f34`)

## v0.2.3 (*2025-02-07*)

✨ **Enhancements**
  - **added trustedProxies settings to API server** (`d3ed25b`)

## v0.2.2 (*2025-01-24*)

📝 **Other Changes**
  - **Change resource ID to actionID, and Workflow action to targetActionID** (`94a46c7`)

## v0.2.1 (*2025-01-24*)

📦 **Updates**
  - **bump pkl to 0.27.2** (`e0683cb`)

📝 **Other Changes**
  - **changed timeoutSeconds -> timeoutDuration** (`f8ed8ed`)

## v0.2.0 (*2025-01-23*)

📦 **Updates**
  - **updated .gitattributes** (`e9ea189`)

📝 **Other Changes**
  - **Use uniform naming convention for {Http,Api,Id,Json} -> {HTTP,API,ID,JSON}** (`d8841da`)

## v0.1.46 (*2025-01-12*)

📝 **Other Changes**
  - **reprioritize request skip validations** (`672a02a`)

## v0.1.45 (*2025-01-11*)

📝 **Other Changes**
  - **removed all deprecated imports** (`7cd9fbe`)

## v0.1.44 (*2025-01-11*)

📝 **Other Changes**
  - **skip & validation is now a listing of any types** (`89cf401`)

## v0.1.43 (*2025-01-11*)

📦 **Updates**
  - **updated README.md** (`809d2c3`)
  - **bump pkl to 0.27.1** (`8183a18`)

## v0.1.42 (*2025-01-11*)

✨ **Enhancements**
  - **upgrade minimum pkl version to 0.26.3. added minor fixes on the relnote generator.** (`c80d5ae`)
  - **Added skip validation helper functions** (`562775e`)

## v0.1.41 (*2025-01-11*)

**📊 Initial Release Statistics:**
- Total commits: 133
- Project inception

**📝 All Changes Since Project Start:**

✨ **Enhancements**
  - **added README.md and relnote generator** (`6c6987c`)
  - **added document pkl module for parsing and creating json, yaml and xml docs** (`ea3a56a`)
  - **added Data resource helper for getting agent data file path** (`efb337e`)
  - **Added file attribute where each associated resource value was saved** (`33eb094`)
  - **added params mapping to http client resources (go gen code)** (`3e93d76`)
  - **added params mapping to http client resources** (`bc81ebb`)
  - **make the build args as mapping; added build env in docker settings** (`ba8c725`)
  - **Add parameters to docker settings; removed unused PKL configurations** (`f8597d4`)
  - **added heroImage and agentIcon** (`008cef9`)
  - **added ollamaImageTag property in workflow settings** (`62fb8d4`)
  - **added condaPackages section** (`be8d611`)
  - **added support for installing conda packages** (`695532e`)
  - **added python script resource** (`d5325c1`)
  - **added install Anaconda option** (`bda10a4`)
  - **added API documentation. Change request API files functions for clarity** (`78b6229`)
  - **added Ubuntu PPA support; add API for querying filetypes by index** (`d45493c`)
  - **added visionFiles to LLM** (`aa2cdf5`)
  - **added API server responseFile** (`c548b8c`)
  - **added request APIs for getting the list of files** (`c717a96`)
  - **added vision and image gen attributes to LLM** (`82340b2`)
  - **Added request file operations. Decode base64 request by default. include allowed HTTP methods.** (`7c53a5d`)
  - **added json response keys to LLM** (`7d9f868`)
  - **added security enforcement settings** (`f0186d5`)
  - **added schema to LLM** (`bdc78d8`)
  - **added extended functions for exec, http and llm modules** (`4bf8fcb`)
  - **added deferred response api** (`fe55044`)
  - **added llm timeout and kdeps dir settings** (`422ea4e`)
  - **add retry mechanism on failed create docs step** (`a85a48d`)
  - **added apiResponse to resource; rename run -> action** (`12e2d39`)
  - **Added API Server Request/Response templates** (`dd71a18`)
  - **added Tags, a globally referenceable token** (`f669c46`)
  - **added hostname and portnum settings to docker** (`df28577`)
  - **Added semantics for external workflows; Additional semantics for resource dependencies** (`d9c1e3f`)
  - **added optional params for default template fields on resource** (`43405b5`)
  - **added docker container init settings** (`5db1d90`)
  - **added dockerGPU and runMode to system config** (`97a4c79`)
  - **updated kdeps template with new fields** (`24495ad`)
  - **added default value for resource env read** (`f4d1453`)
  - **Updated kdeps.pkl to add global configs** (`52cd266`)
  - **Added docker image and llmsettings directive** (`5beaefd`)
  - **added template generation step** (`9644e40`)
  - **try adding prefixed 'v' for go mod tidy** (`74b9fa3`)
  - **added sleep timer to resolve new pkl version release** (`6808149`)
  - **add global settings for kdeps** (`b8da2fd`)
  - **Added API server mode** (`199a52c`)
  - **added project support files** (`a805e3b`)
  - **Added Makefile and gen-go code** (`3acfb33`)
  - **Added initial core schema** (`a61fbf5`)

📦 **Updates**
  - **Revert "upgrade pkl to 0.27.0; register python resource to resource pkl"** (`77fdf1a`)
    This reverts commit 86a334d697479e307513f759d2a7b0b06f9be35c.
  - **Revert "Update Gradle to 8.10.2"** (`95e60f9`)
    This reverts commit 45a93d3ebfa5aa2e795144984f0cde22bc1dc127.
  - **Update Gradle to 8.10.2** (`45a93d3`)
  - **upgrade pkl to 0.27.0; register python resource to resource pkl** (`86a334d`)
  - **upgrade cicd pkl to 0.26.1** (`ea9e7c9`)
  - **updated gen sdk** (`fa77c06`)
  - **updated gen sdks** (`bf7a441`)
  - **updated gen code with dockerGPU and runMode setings** (`2254e19`)
  - **Update CNAME** (`ae8b694`)

📝 **Other Changes**
  - **return stderr if not empty on stdout function** (`4ae7fb9`)
  - **change request function from param("..") -> params("..")** (`541ab82`)
  - **import upstream PKL modules and KDEPS PKL helpers in resource & api response** (`ef98d66`)
  - **document renderers respond with null rather than empty string** (`1d78332`)
  - **decode base64 strings by default on all Resource types** (`80f093c`)
  - **Make the API response errors block listing (array)** (`7283cbe`)
  - **register python resource to resource pkl** (`c13adfa`)
  - **reflect name changes to dockersettings go source** (`4707bd5`)
  - **changed ppa to repositories** (`fc9ae43`)
  - **Make workflow an open module** (`71e4b35`)
  - **deprecate postflightCheck** (`f49439b`)
  - **simplify usage of vision models and attachments** (`a8366ea`)
  - **fix typo for resolving filepath** (`9e69fa3`)
  - **renamed responseFile to file** (`4183cfc`)
  - **renamed responseFile to file** (`5eab2ea`)
  - **always return the first element, in case of a single file upload** (`32eba99`)
  - **support multiple file uploads** (`0e325fa`)
  - **make api serverr request extendable** (`9408e44`)
  - **fix isEmpty on class** (`954d32c`)
  - **simplify resource resolver methods** (`7d9ff35`)
  - **remove api server response type, focus on json response for now** (`bae99a2`)
  - **removed textproto from api response types (gen code)** (`fad23ae`)
  - **removed textproto from api response types** (`6b41ca5`)
  - **response keys gen code** (`06de5e8`)
  - **fix http response block** (`8c1c31e`)
  - **change schema to jsonResponse (bool)** (`4e3c044`)
  - **set empty defaults when calling a function** (`733cfd2`)
  - **increase build timeout to 1 minute** (`56aa39c`)
  - **use a unified api for accessing resource values** (`7789db7`)
  - **exec, chat and client are all blocks** (`77f14c0`)
  - **pre/post flights now a block that includes custom api error** (`797a5b1`)
  - **change resource action to be not a listing** (`f7bb160`)
  - **Revert "disable docs gha action until sorted this out"** (`560985d`)
    This reverts commit 6bed01569525deb1c8fc803b788c9391b1540b01.
  - **disable docs gha action until sorted this out** (`6bed015`)
  - **simplify api server request template by using Dynamic maps** (`c62f4c2`)
  - **workflow and api server has a single required action** (`a974f74`)
  - **non-optional fields and use listing boolean for prefligh & skip steps** (`8e8a242`)
  - **dockerSettings now becomes agentSettings since it has LLM models & hostIP/portNum was moved to apiSettings** (`47bfc0b`)
  - **GEN: move the hostIP and serverPort settings to API Server settings** (`3d8340a`)
  - **move the hostIP and serverPort settings to API Server settings** (`6fcc08e`)
  - **Changed docker hostIP portbinding settings** (`bee501f`)
  - **Changed docker settings for hostName and portNum; Fixed port default value** (`226dc66`)
  - **remove the placeholders for resource and workflow** (`e06e1d7`)
  - **remove the placeholders for resource and workflow** (`8a2fef1`)
  - **Removed llm-apis (for now) and make local LLM the default. (reinstate commercial and cloud llm-apis in future versions)** (`5e46f85`)
  - **Removed workflows array, all pkl files within resources/ folder will be a workflow. Make a resource not an array but a single entry** (`800f898`)
  - **fix workflow template validation** (`646a755`)
  - **removed RAG resource, and expect condition, renamed check to preflight, renamed api to httpclient** (`b1b2b1c`)
  - **removed modelfile, parameters, schema chat** (`902ed91`)
  - **removed interactive input for ENV** (`40f8097`)
  - **renamed ResourceAPI to ResourceHTTPClient** (`3a04936`)
  - **deprecate templates on build** (`ce348d7`)
  - **reinstate read resource** (`48ea4fc`)
  - **init llmapikeys on settings** (`67825b2`)
  - **removed read resource for now** (`8f841d8`)
  - **read llm api keys from env vars** (`281af1f`)
  - **regexp match for go version semantics - PklProject** (`2b77fc9`)
  - **regexp match for go version semantics - GHA** (`f60f8f0`)
  - **regexp match for go version semantics - GHA** (`05a8838`)
  - **regexp match for go version semantics** (`201a555`)
  - **init go.mod; move schema/pkg to schema/gen** (`83c815a`)
  - **initialize the default routes** (`175c9f6`)
  - **renamed Settings.pkl -> Project.pkl** (`720a55f`)
  - **fix indentation** (`80a1347`)
  - **move schema module path to deps/pkl** (`4814b3b`)
  - **specify the source url pattern in Pkl docs** (`9a0864b`)
  - **use pkl-docs env for deploying pages** (`22c31db`)
  - **Delete CNAME** (`4c308c4`)
  - **Create CNAME** (`9d4b2d9`)
  - **Fix resolution of pkl files** (`e3d1588`)
  - **Fix zip package resolution in PklProject; Use pkg as gen folder** (`5355efe`)
  - **fix version resolution** (`c1401f0`)
  - **Initial commit** (`26d33e0`)

---
*Generated on 2025-07-04 06:16:34 by [Enhanced Release Notes Generator](scripts/generate_release_notes.sh)*
