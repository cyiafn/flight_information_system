function fakeStringify(value) {
    let result = "";
    switch (typeof value) {
      case "string":
        result += value;
        break;
      case "number":
        result += value.toString();
        break;
      case "boolean":
        result += value ? "true" : "false";
        break;
      case "object":
        if (value === null) {
          result += "null";
        } else if (Array.isArray(value)) {
          result += "[";
          for (let i = 0; i < value.length; i++) {
            result += fakeStringify(value[i]);
            if (i < value.length - 1) {
              result += ",";
            }
          }
          result += "]";
        } else {
          result += "{";
          let keys = Object.keys(value);
          for (let i = 0; i < keys.length; i++) {
            let key = keys[i];
            let val = value[key];
            result += `"${key}":${fakeStringify(val)}`;
            if (i < keys.length - 1) {
              result += ",";
            }
          }
          result += "}";
        }
        break;
      case "undefined":
        result += "undefined";
        break;
    }
    
    return result;
}

module.exports = fakeStringify;