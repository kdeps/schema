plugins {
  id("org.pkl-lang") version "0.26.1"
}

val maybeVersion = System.getenv("VERSION")

pkl {
  project {
    packagers {
      register("makePackages") {
        if (maybeVersion != null) {
          environmentVariables.put("VERSION", maybeVersion)
        }
        projectDirectories.from(file("deps/pkl/"))
      }
    }
  }
  // ./gradlew pkldoc
  if (maybeVersion != null) {
    pkldocGenerators {
      register("pkldoc") {
        sourceModules =
          listOf(uri("package://schema.kdeps.com/core@$maybeVersion"))
      }
    }
  }
}
