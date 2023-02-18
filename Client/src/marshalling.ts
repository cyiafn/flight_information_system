


export function marshall(data:unknown):any {
  // Store the results as a String
  let strRes = ""; 
  // Crate a list to store the attibutes
  let attrRes = [];


  switch (typeof data) { //Check the data types...
    case "string":
      
      strRes += "" + data + ""

     attrRes.push("0");

      break;
    case "number":
      
      strRes += data.toString();
     attrRes.push("1");

      break;
    case "boolean":
      
      strRes += data ? "true" : "false";
     attrRes.push("2");

      break;
    case "object":

      if (data === null) {
        strRes += "null";
      } 
      // If the data is array
      else if (Array.isArray(data)) {
       attrRes.push("3");
        strRes += data.length;
        strRes += "[";
        for (let i = 0; i < data.length; i++) {
          const {str, attr} = marshall(data[i]);
          strRes += str;
          attrRes.push(...attr);
          if (i < data.length - 1) {
            strRes += ",";
          }
        }
        strRes += "]";
      } 
      // If the data is Object
      else {
       attrRes.push("4");
        strRes += "{";
        let keys = Object.keys(data);
        for(const key of keys) {
          strRes += key + ':'
          const {str, attr} = marshall(data[key as keyof typeof data]);
          strRes += str;
          attrRes.push(...attr);
        }
        strRes += "}";
      }

      break;
    case "undefined":

      strRes += "undefined";
      
      break;
  }
  return { str: strRes, attr : attrRes as string[] }
}