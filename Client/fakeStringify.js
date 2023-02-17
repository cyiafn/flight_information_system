function fakeStringify(value) {
    let result = ""; // Store the results as a String
    let attributes = ""; //Create a list to store the attributes.

    // Stringify will result in attribute-value attribute-value -> 
    // TO-DO Come up with another attribute1,attribute2,value,value
    
    switch (typeof value) { //Check the value types...
      case "string":
       
        result+= '"'
        result += value;
        result += '"';
      
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
    

    
    return (result);
}



module.exports = fakeStringify;