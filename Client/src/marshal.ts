export function marshal(data:unknown):any {
  // Store the results as a String
  let strRes = ""; 
  
  // Check the data types...
  switch (typeof data) { 
    case "string":

      strRes += "" + data + ""

      break;
    case "number":
      
      strRes += data.toString();

      break;
    case "boolean":
      
      strRes += data ? "true" : "false";

      break;
    case "object":

      if (data === null) {
        strRes += "null";
      } 

      // If the data is array
      else if (Array.isArray(data)) {
        strRes += data.length;
        strRes += "[";
        for (let i = 0; i < data.length; i++) {
          const {str} = marshal(data[i]);
          strRes += str;
         
          if (i < data.length - 1)
            strRes += ",";
        }
        strRes += "]";
      } 

      // If the data is Object
      else {
        strRes += "{";
        let keys = Object.keys(data);
        for(const key of keys) {
          strRes += key + ':'
          const {str} = marshal(data[key as keyof typeof data]);
          strRes += str;
        }
        strRes += "}";
      }

      break;
    case "undefined":

      strRes += "undefined";
      
      break;
  }
  return strRes
}