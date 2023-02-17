
function retrieveDataTypes(variable) {
    
  const data = JSON.parse(variable); // Not sure if can parse but i got lazy to find alternative...
  const types = [];

  function getType(item) {
    const type = typeof item;
    if (Array.isArray(item)) { // Check if the item is Array then check its subsequent items
      types.push("array");
      item.forEach((element) => getType(element));
    } else if (type === "object") { // Check if the item is Object then check its subsequent items
      types.push("object");
      Object.values(item).forEach((val) => getType(val));
    } else { // Check normally...
      types.push(type);
    }
  }

  getType(data);
  return types.toString();
}
  

  module.exports = retrieveDataTypes;