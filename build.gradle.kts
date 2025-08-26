plugins {
  id("org.pkl-lang") version "0.29.0"
}

val maybeVersion = System.getenv("VERSION")

pkl {
  project {
    packagers {
      register("makePackages") {
        if (maybeVersion != null) {
          environmentVariables.put("VERSION", maybeVersion)
        }
        projectDirectories.from(
          file("deps/pkl/"),
          file("deps/pkl/external/pkl-pantry/packages/k8s.contrib/"),
          file("deps/pkl/external/pkl-pantry/packages/pkl.pipe/"),
          file("deps/pkl/external/pkl-pantry/packages/org.openapis.v3/"),
          file("deps/pkl/external/pkl-pantry/packages/org.openapis.v3.contrib/")
        )
      }
    }
  }
  // ./gradlew pkldoc
  if (maybeVersion != null) {
    pkldocGenerators {
      register("pkldoc") {
        sourceModules = listOf(uri("package://schema.kdeps.com/core@$maybeVersion"))
      }
    }
  }
}
