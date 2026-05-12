package com.example;

import org.apache.commons.lang3.StringUtils;
import com.google.gson.Gson;
import com.google.gson.JsonObject;

public class Main {
    public static void main(String[] args) {
        String msg = StringUtils.capitalize("hello from the maven uberjar!");
        
        JsonObject json = new JsonObject();
        json.addProperty("message", msg);
        json.addProperty("success", true);
        
        Gson gson = new Gson();
        System.out.println(gson.toJson(json));
    }
}
