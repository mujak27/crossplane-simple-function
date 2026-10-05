import sys
from pathlib import Path
import yaml

crd = yaml.safe_load(Path(sys.argv[1]).read_text())
spec = crd["spec"]

versions = []
for version in spec["versions"]:
    schema = version["schema"]["openAPIV3Schema"]

    # Keep custom fields; Crossplane supplies the XR metadata.
    schema["properties"] = {
        key: value
        for key, value in schema["properties"].items()
        if key in ("spec", "status")
    }

    if "required" in schema:
        schema["required"] = [
            key for key in schema["required"]
            if key in schema["properties"]
        ]
        if not schema["required"]:
            del schema["required"]

    versions.append({
        "name": version["name"],
        "served": version["served"],
        "referenceable": version["storage"],
        "schema": {"openAPIV3Schema": schema},
    })

xrd = {
    "apiVersion": "apiextensions.crossplane.io/v2",
    "kind": "CompositeResourceDefinition",
    "metadata": {"name": crd["metadata"]["name"]},
    "spec": {
        "scope": spec["scope"],
        "group": spec["group"],
        "names": {
            key: spec["names"][key]
            for key in ("kind", "plural", "singular")
            if key in spec["names"]
        },
        "versions": versions,
    },
}

Path(sys.argv[2]).write_text(
    yaml.safe_dump(xrd, sort_keys=False)
)