function deploymentStart(parent, childNode) {
  return new Promise((resolve) => {
    parent.appendChild(childNode);
    animateDeployment(childNode, resolve, "0px");
  });
}

function deploymentEnd(_, childNode) {
  return new Promise((resolve) => 
    animateDeployment(childNode, resolve, "-550px")
  );
}

function animateDeployment(node, resolve, right) {
  $(`div[id='${node.id}']`).animate(
    { right }, 225,  "swing",
    () => setTimeout(resolve, 1_000),
  );
}