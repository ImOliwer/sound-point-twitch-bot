import styled from "styled-components";

type AdditionContainerProps = {
  buttonBackground: string
}

export const Container = styled.div`
  width: 100%;
  height: 100vh;
  background: #26262a;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
`;

export const Title = styled.h2`
  font-family: Arial, sans-serif;
  font-weight: bold;
  font-size: 36px;
  color: #f0f0f0;
  margin-bottom: 50px;
`;

export const AdditionContainer = styled.div<AdditionContainerProps>`
  width: 22.5%;
  height: 30%;
  background: #5f54b9;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;

  & > label {
    font-family: Arial, sans-serif;
    font-size: 16px;
    color: #F0F0F0;
    padding-bottom: 3.5px;
    cursor: pointer;
    border-bottom: 2px solid #fff;
    margin-bottom: 15px;
  }

  & > button {
    padding: 10px;
    width: 60%;
    border: none;
    outline: none;
    color: #fff;
    border-radius: 5px;
    text-transform: uppercase;
    font-family: Arial;
    font-weight: bold;
    background: ${props => props.buttonBackground || "#eb5f5f"};
    margin-top: 15px;
    cursor: pointer;
  }

  & > input {
    outline: none;
    border: none;
    padding: 10px;
    width: 50%;
  }
`;

export const InfoContainer = styled.div`
  width: 60%;
  display: flex;
  flex-direction: column;

  & > input {
    margin-bottom: 5px;
  }
`;

export const CooldownInputContainer = styled.div`
  width: 100%;
  display: flex;

  & > input {
    width: 55%;
    margin-right: 5px;
  }

  & > select {
    flex: 1;
    border-radius: 5px;
    outline: none;
    border: none;
  }
`;